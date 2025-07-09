package aggregator

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Layr-Labs/hourglass-avs-template/internal/consensus"
	"github.com/Layr-Labs/hourglass-avs-template/internal/types"

	log "github.com/sirupsen/logrus"
)

type Aggregator struct {
	minOperators       int
	responseTimeout    time.Duration
	consensusThreshold float64
	consensusEngine    *consensus.ConsensusEngine
	taskStates         map[string]*types.TaskState
	mu                 sync.RWMutex
	coalescedRequests  map[string][]*types.TemperatureTask
	coalesceMu         sync.RWMutex
}

func NewAggregator(minOperators int, responseTimeout time.Duration, consensusThreshold float64, consensusEngine *consensus.ConsensusEngine) *Aggregator {
	return &Aggregator{
		minOperators:       minOperators,
		responseTimeout:    responseTimeout,
		consensusThreshold: consensusThreshold,
		consensusEngine:    consensusEngine,
		taskStates:         make(map[string]*types.TaskState),
		coalescedRequests:  make(map[string][]*types.TemperatureTask),
	}
}

func (a *Aggregator) CreateTask(task types.TemperatureTask) (*types.TaskState, error) {
	taskState := &types.TaskState{
		Task:      task,
		Status:    types.TaskStatusPending,
		Operators: []string{},
		Responses: []types.OperatorResponse{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	a.mu.Lock()
	a.taskStates[task.TaskID] = taskState
	a.mu.Unlock()
	
	log.Infof("Created task %s for location %s (%.2f, %.2f) with threshold %.2f", 
		task.TaskID, task.Location.City, task.Location.Latitude, task.Location.Longitude, task.Threshold)
	
	return taskState, nil
}

func (a *Aggregator) DistributeTask(ctx context.Context, taskID string, availableOperators []string, availableAPIs []string) error {
	a.mu.Lock()
	taskState, exists := a.taskStates[taskID]
	if !exists {
		a.mu.Unlock()
		return fmt.Errorf("task %s not found", taskID)
	}
	
	if taskState.Status != types.TaskStatusPending {
		a.mu.Unlock()
		return fmt.Errorf("task %s is not in pending state", taskID)
	}
	
	selectedOperators := a.selectOperators(availableOperators, taskState.Task.TaskID)
	if len(selectedOperators) < a.minOperators {
		a.mu.Unlock()
		return fmt.Errorf("insufficient operators: %d < %d", len(selectedOperators), a.minOperators)
	}
	
	taskState.Operators = selectedOperators
	taskState.Status = types.TaskStatusDistributed
	taskState.UpdatedAt = time.Now()
	a.mu.Unlock()
	
	_ = a.createTaskDistributions(taskState.Task, selectedOperators, availableAPIs)
	
	log.Infof("Distributing task %s to %d operators", taskID, len(selectedOperators))
	
	return nil
}

func (a *Aggregator) selectOperators(availableOperators []string, taskID string) []string {
	seed := a.generateSeed(taskID)
	rng := rand.New(rand.NewSource(seed))
	
	shuffled := make([]string, len(availableOperators))
	copy(shuffled, availableOperators)
	
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	
	numToSelect := a.minOperators
	if numToSelect > len(shuffled) {
		numToSelect = len(shuffled)
	}
	
	return shuffled[:numToSelect]
}

func (a *Aggregator) generateSeed(taskID string) int64 {
	h := sha256.New()
	h.Write([]byte(taskID))
	hash := h.Sum(nil)
	
	seed := int64(0)
	for i := 0; i < 8 && i < len(hash); i++ {
		seed = (seed << 8) | int64(hash[i])
	}
	
	return seed
}

func (a *Aggregator) createTaskDistributions(task types.TemperatureTask, operators []string, availableAPIs []string) []types.TaskDistribution {
	distributions := make([]types.TaskDistribution, len(operators))
	
	apisPerOperator := len(availableAPIs) / len(operators)
	if apisPerOperator < 1 {
		apisPerOperator = 1
	}
	
	apiIndex := 0
	for i, operator := range operators {
		assignedAPIs := make([]string, 0, apisPerOperator)
		
		for j := 0; j < apisPerOperator && apiIndex < len(availableAPIs); j++ {
			assignedAPIs = append(assignedAPIs, availableAPIs[apiIndex])
			apiIndex++
		}
		
		if apiIndex >= len(availableAPIs) {
			apiIndex = 0
		}
		
		distributions[i] = types.TaskDistribution{
			TaskID:       task.TaskID,
			Task:         task,
			AssignedAPIs: assignedAPIs,
			Deadline:     time.Now().Add(a.responseTimeout),
		}
		
		log.Debugf("Operator %s assigned APIs: %v", operator, assignedAPIs)
	}
	
	return distributions
}

func (a *Aggregator) CollectResponses(ctx context.Context, taskID string, response types.OperatorResponse) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	taskState, exists := a.taskStates[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}
	
	if taskState.Status != types.TaskStatusDistributed && taskState.Status != types.TaskStatusExecuting {
		return fmt.Errorf("task %s is not in correct state for collecting responses", taskID)
	}
	
	if taskState.Status == types.TaskStatusDistributed {
		taskState.Status = types.TaskStatusExecuting
	}
	
	isValidOperator := false
	for _, op := range taskState.Operators {
		if op == response.OperatorID {
			isValidOperator = true
			break
		}
	}
	
	if !isValidOperator {
		return fmt.Errorf("operator %s not assigned to task %s", response.OperatorID, taskID)
	}
	
	for _, existing := range taskState.Responses {
		if existing.OperatorID == response.OperatorID {
			return fmt.Errorf("duplicate response from operator %s", response.OperatorID)
		}
	}
	
	taskState.Responses = append(taskState.Responses, response)
	taskState.UpdatedAt = time.Now()
	
	log.Infof("Collected response from operator %s for task %s (%d/%d responses)", 
		response.OperatorID, taskID, len(taskState.Responses), len(taskState.Operators))
	
	if len(taskState.Responses) >= a.minOperators {
		go a.tryAggregate(taskID)
	}
	
	return nil
}

func (a *Aggregator) tryAggregate(taskID string) {
	time.Sleep(2 * time.Second)
	
	a.mu.Lock()
	taskState, exists := a.taskStates[taskID]
	if !exists || taskState.Status != types.TaskStatusExecuting {
		a.mu.Unlock()
		return
	}
	
	if len(taskState.Responses) < a.minOperators {
		a.mu.Unlock()
		return
	}
	
	taskState.Status = types.TaskStatusAggregating
	a.mu.Unlock()
	
	result, err := a.aggregate(taskState)
	if err != nil {
		log.Errorf("Failed to aggregate task %s: %v", taskID, err)
		a.mu.Lock()
		taskState.Status = types.TaskStatusFailed
		a.mu.Unlock()
		return
	}
	
	a.mu.Lock()
	taskState.ConsensusResult = result
	taskState.Status = types.TaskStatusCompleted
	taskState.UpdatedAt = time.Now()
	a.mu.Unlock()
	
	log.Infof("Task %s completed: %.2f°C, meets threshold: %v, confidence: %.2f", 
		taskID, result.Temperature, result.MeetsThreshold, result.Confidence)
}

func (a *Aggregator) aggregate(taskState *types.TaskState) (*types.ConsensusResult, error) {
	allDataPoints := make([]types.DataPoint, 0)
	
	for _, response := range taskState.Responses {
		for _, dp := range response.DataPoints {
			if a.verifyDataPoint(dp, response.OperatorID, taskState.Task.TaskID) {
				allDataPoints = append(allDataPoints, dp)
			}
		}
	}
	
	if len(allDataPoints) < a.consensusEngine.MinSources {
		return nil, fmt.Errorf("insufficient valid data points: %d", len(allDataPoints))
	}
	
	result, err := a.consensusEngine.ReachConsensus(allDataPoints)
	if err != nil {
		return nil, err
	}
	
	result.TaskID = taskState.Task.TaskID
	a.consensusEngine.VerifyThreshold(result, taskState.Task.Threshold)
	
	result.AggregatedSig = a.aggregateSignatures(taskState.Responses)
	
	return result, nil
}

func (a *Aggregator) verifyDataPoint(dp types.DataPoint, operatorID, taskID string) bool {
	return consensus.VerifySignature(dp.Signature, operatorID, taskID, dp.Temperature)
}

func (a *Aggregator) aggregateSignatures(responses []types.OperatorResponse) []byte {
	h := sha256.New()
	
	for _, resp := range responses {
		h.Write([]byte(resp.OperatorID))
		h.Write(resp.Signature)
		h.Write([]byte(resp.Timestamp.String()))
	}
	
	return h.Sum(nil)
}

func (a *Aggregator) GetTaskState(taskID string) (*types.TaskState, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	taskState, exists := a.taskStates[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}
	
	return taskState, nil
}

func (a *Aggregator) CoalesceRequest(location types.Location) string {
	a.coalesceMu.Lock()
	defer a.coalesceMu.Unlock()
	
	key := fmt.Sprintf("%.6f,%.6f", location.Latitude, location.Longitude)
	
	for existingKey := range a.coalescedRequests {
		if existingKey == key {
			return existingKey
		}
	}
	
	return ""
}

func (a *Aggregator) AddToCoalescedRequest(key string, task *types.TemperatureTask) {
	a.coalesceMu.Lock()
	defer a.coalesceMu.Unlock()
	
	a.coalescedRequests[key] = append(a.coalescedRequests[key], task)
}

func (a *Aggregator) ProcessCoalescedRequests() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		a.coalesceMu.Lock()
		requests := a.coalescedRequests
		a.coalescedRequests = make(map[string][]*types.TemperatureTask)
		a.coalesceMu.Unlock()
		
		for key, tasks := range requests {
			if len(tasks) > 0 {
				log.Infof("Processing %d coalesced requests for location %s", len(tasks), key)
			}
		}
	}
}

func (a *Aggregator) WaitForCompletion(ctx context.Context, taskID string) (*types.ConsensusResult, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	deadline := time.Now().Add(a.responseTimeout)
	
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			a.mu.RLock()
			taskState, exists := a.taskStates[taskID]
			if !exists {
				a.mu.RUnlock()
				return nil, fmt.Errorf("task %s not found", taskID)
			}
			
			if taskState.Status == types.TaskStatusCompleted {
				result := taskState.ConsensusResult
				a.mu.RUnlock()
				return result, nil
			}
			
			if taskState.Status == types.TaskStatusFailed {
				a.mu.RUnlock()
				return nil, fmt.Errorf("task %s failed", taskID)
			}
			
			a.mu.RUnlock()
			
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("timeout waiting for task %s", taskID)
			}
		}
	}
}

type AggregatorMetrics struct {
	TasksCreated     int64
	TasksCompleted   int64
	TasksFailed      int64
	ResponsesCollected int64
	mu               sync.RWMutex
}

func (m *AggregatorMetrics) RecordTask(status types.TaskStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.TasksCreated++
	switch status {
	case types.TaskStatusCompleted:
		m.TasksCompleted++
	case types.TaskStatusFailed:
		m.TasksFailed++
	}
}

func (m *AggregatorMetrics) RecordResponse() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponsesCollected++
}

func (m *AggregatorMetrics) GetStats() (created, completed, failed, responses int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.TasksCreated, m.TasksCompleted, m.TasksFailed, m.ResponsesCollected
}