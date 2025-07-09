package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/hourglass-avs-template/internal/consensus"
	"github.com/Layr-Labs/hourglass-avs-template/internal/datasources"
	"github.com/Layr-Labs/hourglass-avs-template/internal/types"

	log "github.com/sirupsen/logrus"
)

type Executor struct {
	OperatorID    string
	dataManager   *datasources.DataSourceManager
	taskTimeout   time.Duration
	maxConcurrent int
	semaphore     chan struct{}
}

func NewExecutor(operatorID string, dataManager *datasources.DataSourceManager, taskTimeout time.Duration, maxConcurrent int) *Executor {
	return &Executor{
		OperatorID:    operatorID,
		dataManager:   dataManager,
		taskTimeout:   taskTimeout,
		maxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
	}
}

func (e *Executor) ExecuteTask(ctx context.Context, task types.TaskDistribution) (*types.OperatorResponse, error) {
	taskCtx, cancel := context.WithTimeout(ctx, e.taskTimeout)
	defer cancel()
	
	log.Infof("Executor %s starting task %s with %d assigned APIs", 
		e.OperatorID, task.TaskID, len(task.AssignedAPIs))
	
	dataPoints, err := e.fetchTemperatureData(taskCtx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch temperature data: %w", err)
	}
	
	if len(dataPoints) == 0 {
		return nil, fmt.Errorf("no valid data points collected")
	}
	
	response := &types.OperatorResponse{
		OperatorID: e.OperatorID,
		TaskID:     task.TaskID,
		DataPoints: dataPoints,
		Timestamp:  time.Now(),
	}
	
	response.Signature = e.signResponse(response)
	
	log.Infof("Executor %s completed task %s with %d data points", 
		e.OperatorID, task.TaskID, len(dataPoints))
	
	return response, nil
}

func (e *Executor) fetchTemperatureData(ctx context.Context, task types.TaskDistribution) ([]types.DataPoint, error) {
	var (
		dataPoints []types.DataPoint
		mu         sync.Mutex
		wg         sync.WaitGroup
	)
	
	errChan := make(chan error, len(task.AssignedAPIs))
	
	for _, apiName := range task.AssignedAPIs {
		source, ok := e.dataManager.GetSource(apiName)
		if !ok {
			log.Warnf("Unknown data source: %s", apiName)
			continue
		}
		
		wg.Add(1)
		go func(src datasources.WeatherDataSource, name string) {
			defer wg.Done()
			
			select {
			case e.semaphore <- struct{}{}:
				defer func() { <-e.semaphore }()
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}
			
			startTime := time.Now()
			weatherData, err := src.GetTemperature(ctx, task.Task.Location)
			duration := time.Since(startTime)
			
			if err != nil {
				log.Errorf("Failed to fetch from %s: %v (duration: %v)", name, err, duration)
				errChan <- fmt.Errorf("%s: %w", name, err)
				return
			}
			
			dataPoint := types.DataPoint{
				Source:      name,
				Temperature: weatherData.Temperature,
				Timestamp:   weatherData.Timestamp,
				Confidence:  e.calculateConfidence(weatherData, duration),
			}
			
			dataPoint.Signature = consensus.GenerateSignature(e.OperatorID, task.TaskID, dataPoint.Temperature)
			
			mu.Lock()
			dataPoints = append(dataPoints, dataPoint)
			mu.Unlock()
			
			log.Debugf("Fetched from %s: %.2f°C (duration: %v, confidence: %.2f)", 
				name, weatherData.Temperature, duration, dataPoint.Confidence)
		}(source, apiName)
	}
	
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		close(errChan)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}
	
	if len(dataPoints) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("all sources failed: %v", errors)
	}
	
	return dataPoints, nil
}

func (e *Executor) calculateConfidence(data *types.WeatherResponse, fetchDuration time.Duration) float64 {
	confidence := 1.0
	
	if fetchDuration > 5*time.Second {
		confidence *= 0.9
	}
	if fetchDuration > 10*time.Second {
		confidence *= 0.8
	}
	
	age := time.Since(data.Timestamp)
	if age > 5*time.Minute {
		confidence *= 0.9
	}
	if age > 10*time.Minute {
		confidence *= 0.7
	}
	
	return confidence
}

func (e *Executor) signResponse(response *types.OperatorResponse) []byte {
	avgTemp := 0.0
	for _, dp := range response.DataPoints {
		avgTemp += dp.Temperature
	}
	avgTemp /= float64(len(response.DataPoints))
	
	return consensus.GenerateSignature(e.OperatorID, response.TaskID, avgTemp)
}

type ExecutorPool struct {
	executors map[string]*Executor
	mu        sync.RWMutex
}

func NewExecutorPool() *ExecutorPool {
	return &ExecutorPool{
		executors: make(map[string]*Executor),
	}
}

func (p *ExecutorPool) AddExecutor(executor *Executor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.executors[executor.OperatorID] = executor
}

func (p *ExecutorPool) GetExecutor(operatorID string) (*Executor, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	executor, ok := p.executors[operatorID]
	return executor, ok
}

func (p *ExecutorPool) ExecuteTasksConcurrently(ctx context.Context, distributions []types.TaskDistribution) ([]*types.OperatorResponse, error) {
	var (
		responses []*types.OperatorResponse
		mu        sync.Mutex
		wg        sync.WaitGroup
	)
	
	errChan := make(chan error, len(distributions))
	
	for _, dist := range distributions {
		executor, ok := p.GetExecutor(dist.Task.TaskID)
		if !ok {
			continue
		}
		
		wg.Add(1)
		go func(exec *Executor, taskDist types.TaskDistribution) {
			defer wg.Done()
			
			response, err := exec.ExecuteTask(ctx, taskDist)
			if err != nil {
				errChan <- fmt.Errorf("executor %s: %w", exec.OperatorID, err)
				return
			}
			
			mu.Lock()
			responses = append(responses, response)
			mu.Unlock()
		}(executor, dist)
	}
	
	wg.Wait()
	close(errChan)
	
	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}
	
	if len(errors) > 0 {
		log.Warnf("Some executors failed: %v", errors)
	}
	
	return responses, nil
}

type ExecutorMetrics struct {
	TasksExecuted   int64
	TasksSucceeded  int64
	TasksFailed     int64
	TotalFetchTime  time.Duration
	AverageFetchTime time.Duration
	mu              sync.RWMutex
}

func (m *ExecutorMetrics) RecordTask(success bool, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.TasksExecuted++
	if success {
		m.TasksSucceeded++
	} else {
		m.TasksFailed++
	}
	
	m.TotalFetchTime += duration
	m.AverageFetchTime = m.TotalFetchTime / time.Duration(m.TasksExecuted)
}

func (m *ExecutorMetrics) GetStats() (executed, succeeded, failed int64, avgTime time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.TasksExecuted, m.TasksSucceeded, m.TasksFailed, m.AverageFetchTime
}