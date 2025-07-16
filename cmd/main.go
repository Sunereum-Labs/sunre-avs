package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Layr-Labs/hourglass-monorepo/ponos/pkg/performer/server"
	performerV1 "github.com/Layr-Labs/protocol-apis/gen/protos/eigenlayer/hourglass/v1/performer"
	"go.uber.org/zap"

	"github.com/Layr-Labs/hourglass-avs-template/internal/aggregator"
	"github.com/Layr-Labs/hourglass-avs-template/internal/consensus"
	"github.com/Layr-Labs/hourglass-avs-template/internal/datasources"
	"github.com/Layr-Labs/hourglass-avs-template/internal/executor"
	"github.com/Layr-Labs/hourglass-avs-template/internal/insurance"
	"github.com/Layr-Labs/hourglass-avs-template/internal/types"
)

type TaskWorker struct {
	logger          *zap.Logger
	oracle          *WeatherOracle
	claimsProcessor *insurance.ClaimsProcessor
	initialized     bool
}

type WeatherOracle struct {
	aggregator      *aggregator.Aggregator
	executorPool    *executor.ExecutorPool
	dataManager     *datasources.DataSourceManager
	consensusEngine *consensus.ConsensusEngine
	taskCounter     int64
	logger          *zap.Logger
}

// Request types that can be handled
type TaskType string

const (
	TaskTypeWeatherCheck TaskType = "weather_check"
	TaskTypeInsuranceClaim TaskType = "insurance_claim"
	TaskTypeLiveWeatherDemo TaskType = "live_weather_demo"
)

type BaseTaskRequest struct {
	Type TaskType `json:"type"`
}

type WeatherCheckRequest struct {
	Type      TaskType       `json:"type"`
	Location  types.Location `json:"location"`
	Threshold float64        `json:"threshold"`
}

type InsuranceClaimTaskRequest struct {
	Type       TaskType                     `json:"type"`
	ClaimRequest types.InsuranceClaimRequest `json:"claim_request"`
	DemoMode     bool                        `json:"demo_mode,omitempty"`
	DemoScenario string                      `json:"demo_scenario,omitempty"`
}

type LiveWeatherDemoRequest struct {
	Type     TaskType       `json:"type"`
	Location types.Location `json:"location"`
}

func NewTaskWorker(logger *zap.Logger) *TaskWorker {
	return &TaskWorker{
		logger: logger,
	}
}

func (tw *TaskWorker) initializeOracle() error {
	if tw.initialized {
		return nil
	}

	config := &types.Config{
		AVS: struct {
			Aggregator struct {
				MinOperators       int           `yaml:"min_operators"`
				ResponseTimeout    time.Duration `yaml:"response_timeout"`
				ConsensusThreshold float64       `yaml:"consensus_threshold"`
			} `yaml:"aggregator"`
		}{
			Aggregator: struct {
				MinOperators       int           `yaml:"min_operators"`
				ResponseTimeout    time.Duration `yaml:"response_timeout"`
				ConsensusThreshold float64       `yaml:"consensus_threshold"`
			}{
				MinOperators:       3,
				ResponseTimeout:    60 * time.Second,
				ConsensusThreshold: 0.67,
			},
		},
		Consensus: struct {
			MinSources   int     `yaml:"min_sources"`
			MADThreshold float64 `yaml:"mad_threshold"`
			CacheTTL     int     `yaml:"cache_ttl"`
		}{
			MinSources:   3,
			MADThreshold: 2.5,
			CacheTTL:     300,
		},
	}

	// Initialize weather APIs
	weatherAPIs := make(map[string]struct {
		BaseURL   string `yaml:"base_url"`
		RateLimit int    `yaml:"rate_limit"`
		APIKey    string `yaml:"api_key,omitempty"`
	})

	weatherAPIs["openmeteo"] = struct {
		BaseURL   string `yaml:"base_url"`
		RateLimit int    `yaml:"rate_limit"`
		APIKey    string `yaml:"api_key,omitempty"`
	}{
		BaseURL:   "https://api.open-meteo.com/v1",
		RateLimit: 60,
		APIKey:    "",
	}

	// Add Tomorrow.io with the provided API key
	weatherAPIs["tomorrowio"] = struct {
		BaseURL   string `yaml:"base_url"`
		RateLimit int    `yaml:"rate_limit"`
		APIKey    string `yaml:"api_key,omitempty"`
	}{
		BaseURL:   "https://api.tomorrow.io/v4",
		RateLimit: 60,
		APIKey:    "8pDrv1hpHeamM4Cq2OWXFgKMYByz9wyY",
	}

	// Add WeatherAPI.com with the provided API key
	weatherAPIs["weatherapi"] = struct {
		BaseURL   string `yaml:"base_url"`
		RateLimit int    `yaml:"rate_limit"`
		APIKey    string `yaml:"api_key,omitempty"`
	}{
		BaseURL:   "https://api.weatherapi.com/v1",
		RateLimit: 60,
		APIKey:    "963477158a4a42f393f194704250907",
	}

	config.WeatherAPIs = weatherAPIs

	oracle, err := NewWeatherOracle(config, tw.logger)
	if err != nil {
		return fmt.Errorf("failed to create oracle: %w", err)
	}

	tw.oracle = oracle
	tw.claimsProcessor = insurance.NewClaimsProcessor(oracle.consensusEngine)
	tw.initialized = true
	return nil
}

func NewWeatherOracle(config *types.Config, logger *zap.Logger) (*WeatherOracle, error) {
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

	// Create simulated operators
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

	return &WeatherOracle{
		aggregator:      agg,
		executorPool:    execPool,
		dataManager:     dataManager,
		consensusEngine: consensusEngine,
		logger:          logger,
	}, nil
}

func (tw *TaskWorker) ValidateTask(t *performerV1.TaskRequest) error {
	tw.logger.Sugar().Infow("Validating task",
		zap.Any("task", t),
	)

	// Initialize oracle if not already done
	if err := tw.initializeOracle(); err != nil {
		return fmt.Errorf("failed to initialize oracle: %w", err)
	}

	// Decode base request to determine type
	var baseReq BaseTaskRequest
	if err := json.Unmarshal(t.Payload, &baseReq); err != nil {
		return fmt.Errorf("invalid task payload: %w", err)
	}

	switch baseReq.Type {
	case TaskTypeWeatherCheck:
		var taskReq WeatherCheckRequest
		if err := json.Unmarshal(t.Payload, &taskReq); err != nil {
			return fmt.Errorf("invalid weather check request: %w", err)
		}
		return tw.validateWeatherCheck(taskReq)

	case TaskTypeInsuranceClaim:
		var taskReq InsuranceClaimTaskRequest
		if err := json.Unmarshal(t.Payload, &taskReq); err != nil {
			return fmt.Errorf("invalid insurance claim request: %w", err)
		}
		return tw.validateInsuranceClaim(taskReq)

	case TaskTypeLiveWeatherDemo:
		var taskReq LiveWeatherDemoRequest
		if err := json.Unmarshal(t.Payload, &taskReq); err != nil {
			return fmt.Errorf("invalid live weather demo request: %w", err)
		}
		return tw.validateWeatherLocation(taskReq.Location)

	default:
		return fmt.Errorf("unknown task type: %s", baseReq.Type)
	}
}

func (tw *TaskWorker) validateWeatherCheck(req WeatherCheckRequest) error {
	if err := tw.validateWeatherLocation(req.Location); err != nil {
		return err
	}
	if req.Threshold < -100 || req.Threshold > 100 {
		return fmt.Errorf("invalid temperature threshold: %f", req.Threshold)
	}
	return nil
}

func (tw *TaskWorker) validateWeatherLocation(location types.Location) error {
	if location.Latitude < -90 || location.Latitude > 90 {
		return fmt.Errorf("invalid latitude: %f", location.Latitude)
	}
	if location.Longitude < -180 || location.Longitude > 180 {
		return fmt.Errorf("invalid longitude: %f", location.Longitude)
	}
	return nil
}

func (tw *TaskWorker) validateInsuranceClaim(req InsuranceClaimTaskRequest) error {
	policy := req.ClaimRequest.Policy
	
	if policy.PolicyID == "" {
		return fmt.Errorf("policy ID is required")
	}
	
	if policy.CoverageAmount <= 0 {
		return fmt.Errorf("invalid coverage amount: %f", policy.CoverageAmount)
	}
	
	if len(policy.Triggers) == 0 {
		return fmt.Errorf("policy must have at least one trigger")
	}
	
	return nil
}

func (tw *TaskWorker) HandleTask(t *performerV1.TaskRequest) (*performerV1.TaskResponse, error) {
	tw.logger.Sugar().Infow("Handling task",
		zap.Any("task", t),
	)

	// Initialize oracle if not already done
	if err := tw.initializeOracle(); err != nil {
		return nil, fmt.Errorf("failed to initialize oracle: %w", err)
	}

	// Decode base request to determine type
	var baseReq BaseTaskRequest
	if err := json.Unmarshal(t.Payload, &baseReq); err != nil {
		return nil, fmt.Errorf("invalid task payload: %w", err)
	}

	switch baseReq.Type {
	case TaskTypeWeatherCheck:
		return tw.handleWeatherCheck(t)
	case TaskTypeInsuranceClaim:
		return tw.handleInsuranceClaim(t)
	case TaskTypeLiveWeatherDemo:
		return tw.handleLiveWeatherDemo(t)
	default:
		return nil, fmt.Errorf("unknown task type: %s", baseReq.Type)
	}
}

func (tw *TaskWorker) handleWeatherCheck(t *performerV1.TaskRequest) (*performerV1.TaskResponse, error) {
	var taskReq WeatherCheckRequest
	if err := json.Unmarshal(t.Payload, &taskReq); err != nil {
		return nil, fmt.Errorf("invalid weather check request: %w", err)
	}

	// Process weather verification
	ctx := context.Background()
	result, err := tw.oracle.ProcessWeatherVerification(ctx, taskReq.Location, taskReq.Threshold)
	if err != nil {
		return nil, fmt.Errorf("weather verification failed: %w", err)
	}

	// Create response
	response := map[string]interface{}{
		"type":            "weather_check_response",
		"temperature":     result.Temperature,
		"meets_threshold": result.MeetsThreshold,
		"confidence":      result.Confidence,
		"data_points":     len(result.DataPoints),
		"timestamp":       time.Now().Unix(),
	}

	resultBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to encode response: %w", err)
	}

	return &performerV1.TaskResponse{
		TaskId: t.TaskId,
		Result: resultBytes,
	}, nil
}

func (tw *TaskWorker) handleInsuranceClaim(t *performerV1.TaskRequest) (*performerV1.TaskResponse, error) {
	var taskReq InsuranceClaimTaskRequest
	if err := json.Unmarshal(t.Payload, &taskReq); err != nil {
		return nil, fmt.Errorf("invalid insurance claim request: %w", err)
	}

	ctx := context.Background()
	policy := taskReq.ClaimRequest.Policy
	claimDate := taskReq.ClaimRequest.ClaimDate

	var weatherData []types.DataPoint
	var err error

	if taskReq.DemoMode {
		// Use demo data for showcase
		tw.logger.Info("Using demo weather data", 
			zap.String("scenario", taskReq.DemoScenario))
		weatherData = insurance.GenerateDemoWeatherData(
			policy.Location, 
			10, // 10 days of data
			taskReq.DemoScenario,
		)
	} else {
		// Fetch real weather data
		result, err := tw.oracle.ProcessWeatherVerification(ctx, policy.Location, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch weather data: %w", err)
		}
		weatherData = result.DataPoints
	}

	// Process the insurance claim
	claimResponse, err := tw.claimsProcessor.ProcessClaim(
		policy,
		weatherData,
		claimDate,
	)
	if err != nil {
		tw.logger.Sugar().Errorw("Claim processing error", "error", err)
	}

	tw.logger.Sugar().Infow("Insurance claim processed",
		"claimId", claimResponse.ClaimID,
		"status", claimResponse.ClaimStatus,
		"payout", claimResponse.PayoutAmount,
		"triggeredPerils", len(claimResponse.TriggeredPerils),
	)

	// Encode response
	resultBytes, err := json.Marshal(claimResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to encode response: %w", err)
	}

	return &performerV1.TaskResponse{
		TaskId: t.TaskId,
		Result: resultBytes,
	}, nil
}

func (o *WeatherOracle) generateTaskID() string {
	o.taskCounter++
	return fmt.Sprintf("task_%d_%d", time.Now().Unix(), o.taskCounter)
}

func (o *WeatherOracle) ProcessWeatherVerification(ctx context.Context, location types.Location, threshold float64) (*types.ConsensusResult, error) {
	task := types.TemperatureTask{
		TaskID:    o.generateTaskID(),
		Location:  location,
		Threshold: threshold,
		Timestamp: time.Now(),
		ChainID:   1,
	}

	o.logger.Sugar().Infow("Starting weather verification",
		"taskId", task.TaskID,
		"location", location.City,
		"lat", location.Latitude,
		"lon", location.Longitude,
		"threshold", threshold,
	)

	// Create task
	_, err := o.aggregator.CreateTask(task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Distribute task
	operators := []string{"op1", "op2", "op3", "op4", "op5"}
	apis := o.dataManager.GetSourceNames()

	err = o.aggregator.DistributeTask(ctx, task.TaskID, operators, apis)
	if err != nil {
		return nil, fmt.Errorf("failed to distribute task: %w", err)
	}

	// Execute tasks in parallel
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
			Deadline:     time.Now().Add(60 * time.Second),
		}

		go func(e *executor.Executor, td types.TaskDistribution) {
			resp, err := e.ExecuteTask(ctx, td)
			if err != nil {
				errorChan <- fmt.Errorf("executor %s: %w", e.OperatorID, err)
				return
			}
			responseChan <- *resp
		}(exec, taskDist)
	}

	// Collect responses
	collectedCount := 0
	timeout := time.After(65 * time.Second)

	for i := 0; i < len(operators); i++ {
		select {
		case resp := <-responseChan:
			if err := o.aggregator.CollectResponses(ctx, task.TaskID, resp); err != nil {
				o.logger.Sugar().Errorw("Failed to collect response", "error", err)
			} else {
				collectedCount++
			}
		case err := <-errorChan:
			o.logger.Sugar().Errorw("Executor error", "error", err)
		case <-timeout:
			o.logger.Sugar().Warnw("Timeout waiting for responses", "collected", collectedCount)
			break
		}
	}

	if collectedCount < 3 {
		return nil, fmt.Errorf("insufficient responses: %d < 3", collectedCount)
	}

	// Wait for aggregation
	result, err := o.aggregator.WaitForCompletion(ctx, task.TaskID)
	if err != nil {
		return nil, fmt.Errorf("failed to complete aggregation: %w", err)
	}

	o.logger.Sugar().Infow("Task completed",
		"taskId", task.TaskID,
		"temperature", result.Temperature,
		"meetsThreshold", result.MeetsThreshold,
		"confidence", result.Confidence,
	)

	return result, nil
}

func (tw *TaskWorker) handleLiveWeatherDemo(t *performerV1.TaskRequest) (*performerV1.TaskResponse, error) {
	var taskReq LiveWeatherDemoRequest
	if err := json.Unmarshal(t.Payload, &taskReq); err != nil {
		return nil, fmt.Errorf("invalid live weather demo request: %w", err)
	}

	// Get current weather data
	ctx := context.Background()
	currentResult, err := tw.oracle.ProcessWeatherVerification(ctx, taskReq.Location, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get current weather: %w", err)
	}

	// Get historical weather data (4 hours ago)
	// For demo purposes, simulate historical data
	historicalTemp := currentResult.Temperature - 2.5 // Assume it was 2.5°C cooler 4 hours ago
	
	// Group data by source
	sourceData := make(map[string]map[string]interface{})
	for _, dp := range currentResult.DataPoints {
		sourceData[dp.Source] = map[string]interface{}{
			"current_temperature": dp.Temperature,
			"current_timestamp": dp.Timestamp,
			"confidence": dp.Confidence,
		}
	}

	// Create a dynamic insurance policy based on current NYC weather
	var policyExample map[string]interface{}
	
	if currentResult.Temperature < 0 {
		// Cold weather - travel insurance
		policyExample = map[string]interface{}{
			"type": "travel_insurance",
			"scenario": "Winter Flight Delay Protection",
			"description": "With current NYC temperature at " + fmt.Sprintf("%.1f°C", currentResult.Temperature) + ", flight delays are more likely",
			"coverage": "$500 per day for weather-related delays",
			"trigger": "Temperature below 0°C causing flight delays",
			"premium": "$25 for 5-day coverage",
			"relevance": "High - Current conditions match policy trigger",
		}
	} else if currentResult.Temperature > 30 {
		// Hot weather - event insurance
		policyExample = map[string]interface{}{
			"type": "event_insurance",
			"scenario": "Outdoor Event Heat Protection",
			"description": "With NYC at " + fmt.Sprintf("%.1f°C", currentResult.Temperature) + ", outdoor events need heat coverage",
			"coverage": "$100,000 event cancellation insurance",
			"trigger": "Temperature above 35°C for event hours",
			"premium": "$2,000 for single event",
			"relevance": "Medium - Close to trigger threshold",
		}
	} else {
		// Moderate weather - property insurance
		policyExample = map[string]interface{}{
			"type": "property_insurance", 
			"scenario": "Weather Damage Protection",
			"description": "Current moderate NYC weather (" + fmt.Sprintf("%.1f°C", currentResult.Temperature) + ") ideal for annual coverage",
			"coverage": "$50,000 weather damage protection",
			"trigger": "Extreme weather events: hail, flooding, wind damage",
			"premium": "$500 annual",
			"relevance": "Standard coverage for unpredictable weather",
		}
	}

	// Create response
	response := map[string]interface{}{
		"type": "live_weather_demo_response",
		"location": map[string]interface{}{
			"city": taskReq.Location.City,
			"latitude": taskReq.Location.Latitude,
			"longitude": taskReq.Location.Longitude,
		},
		"current_weather": map[string]interface{}{
			"temperature": currentResult.Temperature,
			"consensus_confidence": currentResult.Confidence,
			"timestamp": time.Now().Unix(),
			"data_sources": sourceData,
		},
		"historical_weather": map[string]interface{}{
			"temperature_4h_ago": historicalTemp,
			"temperature_change": currentResult.Temperature - historicalTemp,
			"trend": func() string {
				if currentResult.Temperature > historicalTemp {
					return "warming"
				} else if currentResult.Temperature < historicalTemp {
					return "cooling"
				}
				return "stable"
			}(),
		},
		"insurance_recommendation": policyExample,
		"consensus_details": map[string]interface{}{
			"algorithm": "MAD (Median Absolute Deviation)",
			"total_sources": len(currentResult.DataPoints),
			"sources_used": func() []string {
				sources := make([]string, 0)
				for _, dp := range currentResult.DataPoints {
					sources = append(sources, dp.Source)
				}
				return sources
			}(),
		},
	}

	resultBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to encode response: %w", err)
	}

	return &performerV1.TaskResponse{
		TaskId: t.TaskId,
		Result: resultBytes,
	}, nil
}

func (o *WeatherOracle) getAPISubsetForOperator(operatorIndex int, apis []string, numOperators int) []string {
	if len(apis) == 0 {
		return []string{}
	}

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
		if end > len(apis) {
			end = len(apis)
		}
	}

	return apis[start:end]
}

func main() {
	ctx := context.Background()
	l, _ := zap.NewProduction()

	w := NewTaskWorker(l)

	pp, err := server.NewPonosPerformerWithRpcServer(&server.PonosPerformerConfig{
		Port:    8080,
		Timeout: 5 * time.Second,
	}, w, l)
	if err != nil {
		panic(fmt.Errorf("failed to create performer: %w", err))
	}

	l.Info("Starting Weather Insurance AVS Performer", zap.Int("port", 8080))

	if err := pp.Start(ctx); err != nil {
		panic(err)
	}
}