package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"temperature-oracle/internal/aggregator"
	"temperature-oracle/internal/consensus"
	"temperature-oracle/internal/datasources"
	"temperature-oracle/internal/executor"
	"temperature-oracle/internal/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	tasksProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "temperature_oracle_tasks_processed_total",
			Help: "Total number of temperature verification tasks processed",
		},
		[]string{"status"},
	)
	
	taskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "temperature_oracle_task_duration_seconds",
			Help:    "Duration of temperature verification tasks",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"phase"},
	)
	
	consensusTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_oracle_consensus_temperature",
			Help: "Consensus temperature for locations",
		},
		[]string{"location"},
	)
)

func init() {
	prometheus.MustRegister(tasksProcessed, taskDuration, consensusTemperature)
}

type TemperatureOracle struct {
	config          *types.Config
	aggregator      *aggregator.Aggregator
	executorPool    *executor.ExecutorPool
	dataManager     *datasources.DataSourceManager
	consensusEngine *consensus.ConsensusEngine
	taskCounter     int64
	mu              sync.Mutex
}

func NewTemperatureOracle(config *types.Config) (*TemperatureOracle, error) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	
	cacheTTL := time.Duration(config.Consensus.CacheTTL) * time.Second
	dataManager := datasources.NewDataSourceManager(config.WeatherAPIs, cacheTTL)
	
	if len(dataManager.GetAllSources()) == 0 {
		return nil, fmt.Errorf("no weather data sources configured")
	}
	
	consensusEngine := consensus.NewConsensusEngine(
		config.Consensus.MinSources,
		config.Consensus.MADThreshold,
	)
	
	agg := aggregator.NewAggregator(
		config.AVS.Aggregator.MinOperators,
		config.AVS.Aggregator.ResponseTimeout,
		config.AVS.Aggregator.ConsensusThreshold,
		consensusEngine,
	)
	
	execPool := executor.NewExecutorPool()
	
	operators := []string{"op1", "op2", "op3", "op4", "op5"}
	for _, opID := range operators {
		exec := executor.NewExecutor(
			opID,
			dataManager,
			60*time.Second,
			3,
		)
		execPool.AddExecutor(exec)
	}
	
	return &TemperatureOracle{
		config:          config,
		aggregator:      agg,
		executorPool:    execPool,
		dataManager:     dataManager,
		consensusEngine: consensusEngine,
	}, nil
}

func (o *TemperatureOracle) generateTaskID() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.taskCounter++
	return fmt.Sprintf("task_%d_%d", time.Now().Unix(), o.taskCounter)
}

func (o *TemperatureOracle) ProcessTemperatureVerification(ctx context.Context, location types.Location, threshold float64) (*types.ConsensusResult, error) {
	start := time.Now()
	
	task := types.TemperatureTask{
		TaskID:    o.generateTaskID(),
		Location:  location,
		Threshold: threshold,
		Timestamp: time.Now(),
		ChainID:   1,
	}
	
	log.Infof("Starting temperature verification task %s for %s (%.2f, %.2f) with threshold %.2f°C",
		task.TaskID, location.City, location.Latitude, location.Longitude, threshold)
	
	timer := prometheus.NewTimer(taskDuration.WithLabelValues("event_detection"))
	_, err := o.aggregator.CreateTask(task)
	if err != nil {
		tasksProcessed.WithLabelValues("failed").Inc()
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	timer.ObserveDuration()
	
	timer = prometheus.NewTimer(taskDuration.WithLabelValues("distribution"))
	operators := []string{"op1", "op2", "op3", "op4", "op5"}
	apis := o.dataManager.GetSourceNames()
	
	err = o.aggregator.DistributeTask(ctx, task.TaskID, operators, apis)
	if err != nil {
		tasksProcessed.WithLabelValues("failed").Inc()
		return nil, fmt.Errorf("failed to distribute task: %w", err)
	}
	timer.ObserveDuration()
	
	timer = prometheus.NewTimer(taskDuration.WithLabelValues("execution"))
	var wg sync.WaitGroup
	responseChan := make(chan types.OperatorResponse, len(operators))
	errorChan := make(chan error, len(operators))
	
	for i, opID := range operators {
		exec, ok := o.executorPool.GetExecutor(opID)
		if !ok {
			continue
		}
		
		apiSubset := o.getAPISubsetForOperator(i, apis, len(operators))
		taskDist := types.TaskDistribution{
			TaskID:       task.TaskID,
			Task:         task,
			AssignedAPIs: apiSubset,
			Deadline:     time.Now().Add(o.config.AVS.Aggregator.ResponseTimeout),
		}
		
		wg.Add(1)
		go func(e *executor.Executor, td types.TaskDistribution) {
			defer wg.Done()
			
			resp, err := e.ExecuteTask(ctx, td)
			if err != nil {
				errorChan <- fmt.Errorf("executor %s: %w", e.OperatorID, err)
				return
			}
			
			responseChan <- *resp
		}(exec, taskDist)
	}
	
	go func() {
		wg.Wait()
		close(responseChan)
		close(errorChan)
	}()
	
	collectedCount := 0
	for resp := range responseChan {
		if err := o.aggregator.CollectResponses(ctx, task.TaskID, resp); err != nil {
			log.Errorf("Failed to collect response: %v", err)
		} else {
			collectedCount++
		}
	}
	
	for err := range errorChan {
		log.Errorf("Executor error: %v", err)
	}
	
	timer.ObserveDuration()
	
	if collectedCount < o.config.AVS.Aggregator.MinOperators {
		tasksProcessed.WithLabelValues("failed").Inc()
		return nil, fmt.Errorf("insufficient responses: %d < %d", collectedCount, o.config.AVS.Aggregator.MinOperators)
	}
	
	timer = prometheus.NewTimer(taskDuration.WithLabelValues("aggregation"))
	result, err := o.aggregator.WaitForCompletion(ctx, task.TaskID)
	if err != nil {
		tasksProcessed.WithLabelValues("failed").Inc()
		return nil, fmt.Errorf("failed to complete aggregation: %w", err)
	}
	timer.ObserveDuration()
	
	tasksProcessed.WithLabelValues("completed").Inc()
	consensusTemperature.WithLabelValues(location.City).Set(result.Temperature)
	
	elapsed := time.Since(start)
	log.Infof("Task %s completed in %v: Temperature=%.2f°C, MeetsThreshold=%v, Confidence=%.2f",
		task.TaskID, elapsed, result.Temperature, result.MeetsThreshold, result.Confidence)
	
	return result, nil
}

func (o *TemperatureOracle) getAPISubsetForOperator(operatorIndex int, apis []string, numOperators int) []string {
	apisPerOperator := len(apis) / numOperators
	if apisPerOperator < 1 {
		apisPerOperator = 1
	}
	
	start := operatorIndex * apisPerOperator
	end := start + apisPerOperator
	
	if end > len(apis) {
		end = len(apis)
	}
	
	if start >= len(apis) {
		start = 0
		end = apisPerOperator
	}
	
	return apis[start:end]
}

func loadConfig(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	if config.AVS.Aggregator.ResponseTimeout == 0 {
		config.AVS.Aggregator.ResponseTimeout = 60 * time.Second
	}
	
	for name, apiConfig := range config.WeatherAPIs {
		if apiKey := os.Getenv(fmt.Sprintf("%s_API_KEY", strings.ToUpper(name))); apiKey != "" {
			apiConfig.APIKey = apiKey
			config.WeatherAPIs[name] = apiConfig
		}
	}
	
	return &config, nil
}

func parseLocation(locationStr string) (types.Location, error) {
	location := types.Location{}
	
	knownLocations := map[string]types.Location{
		"new york":      {Latitude: 40.7128, Longitude: -74.0060, City: "New York", Country: "USA"},
		"london":        {Latitude: 51.5074, Longitude: -0.1278, City: "London", Country: "UK"},
		"tokyo":         {Latitude: 35.6762, Longitude: 139.6503, City: "Tokyo", Country: "Japan"},
		"paris":         {Latitude: 48.8566, Longitude: 2.3522, City: "Paris", Country: "France"},
		"sydney":        {Latitude: -33.8688, Longitude: 151.2093, City: "Sydney", Country: "Australia"},
		"san francisco": {Latitude: 37.7749, Longitude: -122.4194, City: "San Francisco", Country: "USA"},
		"singapore":     {Latitude: 1.3521, Longitude: 103.8198, City: "Singapore", Country: "Singapore"},
		"dubai":         {Latitude: 25.2048, Longitude: 55.2708, City: "Dubai", Country: "UAE"},
	}
	
	locationLower := strings.ToLower(locationStr)
	if loc, ok := knownLocations[locationLower]; ok {
		return loc, nil
	}
	
	var lat, lon float64
	if _, err := fmt.Sscanf(locationStr, "%f,%f", &lat, &lon); err == nil {
		location.Latitude = lat
		location.Longitude = lon
		location.City = fmt.Sprintf("Location (%.2f, %.2f)", lat, lon)
		location.Country = "Unknown"
		return location, nil
	}
	
	return location, fmt.Errorf("unknown location format: %s", locationStr)
}

func main() {
	var (
		configPath  = flag.String("config", "config/config.yaml", "Path to configuration file")
		location    = flag.String("location", "", "Location (city name or lat,lon)")
		threshold   = flag.Float64("threshold", 25.0, "Temperature threshold in Celsius")
		metricsPort = flag.String("metrics-port", "8080", "Port for Prometheus metrics")
		logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()
	
	if *location == "" {
		log.Fatal("Location is required. Use --location flag")
	}
	
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)
	
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	oracle, err := NewTemperatureOracle(config)
	if err != nil {
		log.Fatalf("Failed to create oracle: %v", err)
	}
	
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		log.Infof("Starting metrics server on :%s", *metricsPort)
		if err := http.ListenAndServe(":"+*metricsPort, nil); err != nil {
			log.Errorf("Metrics server error: %v", err)
		}
	}()
	
	loc, err := parseLocation(*location)
	if err != nil {
		log.Fatalf("Failed to parse location: %v", err)
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		log.Info("Received shutdown signal")
		cancel()
	}()
	
	log.Infof("Temperature Oracle started")
	log.Infof("Location: %s (%.2f, %.2f)", loc.City, loc.Latitude, loc.Longitude)
	log.Infof("Threshold: %.2f°C", *threshold)
	log.Infof("Available data sources: %v", oracle.dataManager.GetSourceNames())
	
	result, err := oracle.ProcessTemperatureVerification(ctx, loc, *threshold)
	if err != nil {
		log.Fatalf("Temperature verification failed: %v", err)
	}
	
	fmt.Printf("\n=== Temperature Verification Result ===\n")
	fmt.Printf("Location: %s (%.2f, %.2f)\n", loc.City, loc.Latitude, loc.Longitude)
	fmt.Printf("Consensus Temperature: %.2f°C\n", result.Temperature)
	fmt.Printf("Threshold: %.2f°C\n", *threshold)
	fmt.Printf("Meets Threshold: %v\n", result.MeetsThreshold)
	fmt.Printf("Confidence: %.2f%%\n", result.Confidence*100)
	fmt.Printf("Data Sources Used: %d\n", len(result.DataPoints))
	fmt.Printf("\nData Points:\n")
	for _, dp := range result.DataPoints {
		fmt.Printf("  - %s: %.2f°C (confidence: %.2f)\n", dp.Source, dp.Temperature, dp.Confidence)
	}
	fmt.Printf("\nTask ID: %s\n", result.TaskID)
	fmt.Printf("Aggregated Signature: %x\n", result.AggregatedSig[:16])
	fmt.Printf("=====================================\n")
	
	log.Info("Temperature oracle shutdown complete")
}