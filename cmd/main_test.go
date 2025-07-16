package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/hourglass-monorepo/ponos/pkg/performer/server"
	performerV1 "github.com/Layr-Labs/protocol-apis/gen/protos/eigenlayer/hourglass/v1/performer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Layr-Labs/hourglass-avs-template/internal/types"
)

// Test the TaskWorker ValidateTask method
func TestValidateTask(t *testing.T) {
	logger := zap.NewNop() // Use no-op logger for tests to reduce noise
	tw := NewTaskWorker(logger)

	tests := []struct {
		name    string
		payload json.RawMessage
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid weather check",
			payload: json.RawMessage(`{
				"type": "weather_check",
				"location": {
					"latitude": 40.7128,
					"longitude": -74.0060,
					"city": "New York",
					"country": "USA"
				},
				"threshold": 25.0
			}`),
			wantErr: false,
		},
		{
			name: "invalid latitude",
			payload: json.RawMessage(`{
				"type": "weather_check",
				"location": {
					"latitude": 91.0,
					"longitude": -74.0060,
					"city": "New York",
					"country": "USA"
				},
				"threshold": 25.0
			}`),
			wantErr: true,
			errMsg:  "invalid latitude",
		},
		{
			name: "invalid longitude",
			payload: json.RawMessage(`{
				"type": "weather_check",
				"location": {
					"latitude": 40.7128,
					"longitude": -181.0,
					"city": "New York",
					"country": "USA"
				},
				"threshold": 25.0
			}`),
			wantErr: true,
			errMsg:  "invalid longitude",
		},
		{
			name: "invalid temperature threshold",
			payload: json.RawMessage(`{
				"type": "weather_check",
				"location": {
					"latitude": 40.7128,
					"longitude": -74.0060,
					"city": "New York",
					"country": "USA"
				},
				"threshold": 150.0
			}`),
			wantErr: true,
			errMsg:  "invalid temperature threshold",
		},
		{
			name: "valid insurance claim",
			payload: json.RawMessage(`{
				"type": "insurance_claim",
				"claim_request": {
					"policy_id": "TEST-001",
					"policy": {
						"policy_id": "TEST-001",
						"policy_holder": "Test Holder",
						"insurance_type": "crop",
						"location": {
							"latitude": 40.7128,
							"longitude": -74.0060,
							"city": "New York",
							"country": "USA"
						},
						"coverage_amount": 100000,
						"premium": 5000,
						"start_date": "2024-01-01T00:00:00Z",
						"end_date": "2024-12-31T00:00:00Z",
						"triggers": [{
							"trigger_id": "HEAT-001",
							"peril": "heat_wave",
							"conditions": {
								"temperature_max": 35
							},
							"payout_ratio": 0.5,
							"description": "Heat wave trigger"
						}]
					},
					"claim_date": "2024-07-15T00:00:00Z",
					"automated_check": true
				},
				"demo_mode": true,
				"demo_scenario": "heat_wave"
			}`),
			wantErr: false,
		},
		{
			name: "missing policy ID",
			payload: json.RawMessage(`{
				"type": "insurance_claim",
				"claim_request": {
					"policy": {
						"policy_id": "",
						"coverage_amount": 100000,
						"triggers": []
					}
				}
			}`),
			wantErr: true,
			errMsg:  "policy ID is required",
		},
		{
			name: "invalid coverage amount",
			payload: json.RawMessage(`{
				"type": "insurance_claim",
				"claim_request": {
					"policy": {
						"policy_id": "TEST-001",
						"coverage_amount": -1000,
						"triggers": [{
							"trigger_id": "TEST",
							"peril": "heat_wave",
							"payout_ratio": 0.5
						}]
					}
				}
			}`),
			wantErr: true,
			errMsg:  "invalid coverage amount",
		},
		{
			name: "no triggers",
			payload: json.RawMessage(`{
				"type": "insurance_claim",
				"claim_request": {
					"policy": {
						"policy_id": "TEST-001",
						"coverage_amount": 100000,
						"triggers": []
					}
				}
			}`),
			wantErr: true,
			errMsg:  "must have at least one trigger",
		},
		{
			name: "valid live weather demo",
			payload: json.RawMessage(`{
				"type": "live_weather_demo",
				"location": {
					"latitude": 40.7128,
					"longitude": -74.0060,
					"city": "New York",
					"country": "USA"
				}
			}`),
			wantErr: false,
		},
		{
			name: "unknown task type",
			payload: json.RawMessage(`{
				"type": "unknown_task"
			}`),
			wantErr: true,
			errMsg:  "unknown task type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &performerV1.TaskRequest{
				TaskId:  []byte("test-task-" + tt.name),
				Payload: tt.payload,
			}

			err := tw.ValidateTask(task)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test the TaskWorker HandleTask method
func TestHandleTask(t *testing.T) {
	logger := zap.NewNop()
	tw := NewTaskWorker(logger)

	tests := []struct {
		name           string
		payload        json.RawMessage
		wantErr        bool
		validateResult func(t *testing.T, result []byte)
	}{
		// Note: Weather check test commented out as it requires external API calls
		// For production tests, use mock data sources or test environment
		// {
		// 	name: "weather check task",
		// 	payload: json.RawMessage(`{
		// 		"type": "weather_check",
		// 		"location": {
		// 			"latitude": 40.7128,
		// 			"longitude": -74.0060,
		// 			"city": "New York",
		// 			"country": "USA"
		// 		},
		// 		"threshold": 25.0
		// 	}`),
		// 	wantErr: false,
		// 	validateResult: func(t *testing.T, result []byte) {
		// 		var response map[string]interface{}
		// 		err := json.Unmarshal(result, &response)
		// 		require.NoError(t, err)
		// 		
		// 		assert.Equal(t, "weather_check_response", response["type"])
		// 		assert.NotNil(t, response["temperature"])
		// 		assert.NotNil(t, response["meets_threshold"])
		// 		assert.NotNil(t, response["confidence"])
		// 		assert.NotNil(t, response["data_points"])
		// 		assert.NotNil(t, response["timestamp"])
		// 	},
		// },
		{
			name: "insurance claim task - heat wave",
			payload: json.RawMessage(`{
				"type": "insurance_claim",
				"claim_request": {
					"policy_id": "TEST-HEAT-001",
					"policy": {
						"policy_id": "TEST-HEAT-001",
						"policy_holder": "Test Farm",
						"insurance_type": "crop",
						"location": {
							"latitude": 40.7128,
							"longitude": -74.0060,
							"city": "New York",
							"country": "USA"
						},
						"coverage_amount": 100000,
						"premium": 5000,
						"start_date": "2024-01-01T00:00:00Z",
						"end_date": "2024-12-31T00:00:00Z",
						"triggers": [{
							"trigger_id": "HEAT-3DAY",
							"peril": "heat_wave",
							"conditions": {
								"temperature_max": 35,
								"consecutive_days": 3
							},
							"payout_ratio": 0.5,
							"description": "Heat wave protection"
						}]
					},
					"claim_date": "2024-07-15T00:00:00Z",
					"automated_check": true
				},
				"demo_mode": true,
				"demo_scenario": "heat_wave"
			}`),
			wantErr: false,
			validateResult: func(t *testing.T, result []byte) {
				var response types.InsuranceClaimResponse
				err := json.Unmarshal(result, &response)
				require.NoError(t, err)
				
				assert.NotEmpty(t, response.ClaimID)
				assert.NotEmpty(t, response.ClaimStatus)
				assert.NotNil(t, response.PayoutAmount)
				assert.NotNil(t, response.TriggeredPerils)
				assert.NotEmpty(t, response.VerificationHash)
			},
		},
		// Note: Live weather demo test is commented out because it requires real API calls
		// which may fail in test environments without proper API keys or network access
		// {
		// 	name: "live weather demo task",
		// 	payload: json.RawMessage(`{
		// 		"type": "live_weather_demo",
		// 		"location": {
		// 			"latitude": 40.7128,
		// 			"longitude": -74.0060,
		// 			"city": "New York",
		// 			"country": "USA"
		// 		}
		// 	}`),
		// 	wantErr: false,
		// 	validateResult: func(t *testing.T, result []byte) {
		// 		var response map[string]interface{}
		// 		err := json.Unmarshal(result, &response)
		// 		require.NoError(t, err)
		// 		
		// 		assert.Equal(t, "live_weather_demo_response", response["type"])
		// 		assert.NotNil(t, response["location"])
		// 		assert.NotNil(t, response["current_weather"])
		// 		assert.NotNil(t, response["historical_weather"])
		// 		assert.NotNil(t, response["insurance_recommendation"])
		// 		assert.NotNil(t, response["consensus_details"])
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &performerV1.TaskRequest{
				TaskId:  []byte("test-task-" + tt.name),
				Payload: tt.payload,
			}

			response, err := tw.HandleTask(task)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if !assert.NoError(t, err) {
					return // Skip further assertions if there's an error
				}
				if !assert.NotNil(t, response) {
					return // Skip further assertions if response is nil
				}
				assert.Equal(t, task.TaskId, response.TaskId)
				assert.NotEmpty(t, response.Result)
				
				if tt.validateResult != nil {
					tt.validateResult(t, response.Result)
				}
			}
		})
	}
}

// Test the weather oracle initialization
func TestWeatherOracleInitialization(t *testing.T) {
	logger := zap.NewNop()
	tw := NewTaskWorker(logger)
	
	// Test initialization
	err := tw.initializeOracle()
	assert.NoError(t, err)
	assert.NotNil(t, tw.oracle)
	assert.NotNil(t, tw.claimsProcessor)
	assert.True(t, tw.initialized)
	
	// Test that second initialization doesn't recreate
	oracle1 := tw.oracle
	err = tw.initializeOracle()
	assert.NoError(t, err)
	assert.Equal(t, oracle1, tw.oracle) // Should be the same instance
}

// Test weather location validation
func TestValidateWeatherLocation(t *testing.T) {
	logger := zap.NewNop()
	tw := NewTaskWorker(logger)

	tests := []struct {
		name     string
		location types.Location
		wantErr  bool
	}{
		{
			name: "valid location",
			location: types.Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
				City:      "New York",
				Country:   "USA",
			},
			wantErr: false,
		},
		{
			name: "latitude too high",
			location: types.Location{
				Latitude:  91.0,
				Longitude: -74.0060,
			},
			wantErr: true,
		},
		{
			name: "latitude too low",
			location: types.Location{
				Latitude:  -91.0,
				Longitude: -74.0060,
			},
			wantErr: true,
		},
		{
			name: "longitude too high",
			location: types.Location{
				Latitude:  40.7128,
				Longitude: 181.0,
			},
			wantErr: true,
		},
		{
			name: "longitude too low",
			location: types.Location{
				Latitude:  40.7128,
				Longitude: -181.0,
			},
			wantErr: true,
		},
		{
			name: "edge case - north pole",
			location: types.Location{
				Latitude:  90.0,
				Longitude: 0.0,
			},
			wantErr: false,
		},
		{
			name: "edge case - south pole",
			location: types.Location{
				Latitude:  -90.0,
				Longitude: 0.0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tw.validateWeatherLocation(tt.location)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test the server integration
func TestServerIntegration(t *testing.T) {
	logger := zap.NewNop()
	tw := NewTaskWorker(logger)

	// Create a test server
	pp, err := server.NewPonosPerformerWithRpcServer(&server.PonosPerformerConfig{
		Port:    8081, // Use different port for testing
		Timeout: 5 * time.Second,
	}, tw, logger)
	require.NoError(t, err)

	// Start the server in a goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		err := pp.Start(ctx)
		if err != nil && err != context.Canceled {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(1 * time.Second)

	// Test that server is ready to handle requests
	// In a real test, you would make HTTP requests to the server
	// For now, we just verify it started without panic
	assert.NotNil(t, pp)
}

// Test concurrent task handling with validation only (no external calls)
func TestConcurrentTaskValidation(t *testing.T) {
	logger := zap.NewNop()
	tw := NewTaskWorker(logger)

	// Create multiple tasks
	numTasks := 50
	tasks := make([]*performerV1.TaskRequest, numTasks)
	for i := 0; i < numTasks; i++ {
		payload := json.RawMessage(`{
			"type": "weather_check",
			"location": {
				"latitude": 40.7128,
				"longitude": -74.0060,
				"city": "New York",
				"country": "USA"
			},
			"threshold": 25.0
		}`)
		
		tasks[i] = &performerV1.TaskRequest{
			TaskId:  []byte(fmt.Sprintf("task-%d", i)),
			Payload: payload,
		}
	}

	// Validate tasks concurrently
	results := make(chan error, len(tasks))

	for _, task := range tasks {
		go func(t *performerV1.TaskRequest) {
			err := tw.ValidateTask(t)
			results <- err
		}(task)
	}

	// Collect results
	successCount := 0
	errorCount := 0
	timeout := time.After(5 * time.Second)

	for i := 0; i < len(tasks); i++ {
		select {
		case err := <-results:
			if err == nil {
				successCount++
			} else {
				errorCount++
			}
		case <-timeout:
			t.Fatal("Timeout waiting for task validation")
		}
	}

	// All tasks should validate successfully
	assert.Equal(t, numTasks, successCount)
	assert.Equal(t, 0, errorCount)
}

// Benchmark task validation
func BenchmarkValidateTask(b *testing.B) {
	logger := zap.NewNop()
	tw := NewTaskWorker(logger)

	payload := json.RawMessage(`{
		"type": "weather_check",
		"location": {
			"latitude": 40.7128,
			"longitude": -74.0060,
			"city": "New York",
			"country": "USA"
		},
		"threshold": 25.0
	}`)

	task := &performerV1.TaskRequest{
		TaskId:  []byte("benchmark-task"),
		Payload: payload,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tw.ValidateTask(task)
	}
}

// Benchmark task handling
func BenchmarkHandleTask(b *testing.B) {
	logger := zap.NewNop()
	tw := NewTaskWorker(logger)

	payload := json.RawMessage(`{
		"type": "weather_check",
		"location": {
			"latitude": 40.7128,
			"longitude": -74.0060,
			"city": "New York",
			"country": "USA"
		},
		"threshold": 25.0
	}`)

	task := &performerV1.TaskRequest{
		TaskId:  []byte("benchmark-task"),
		Payload: payload,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tw.HandleTask(task)
	}
}

// Test_TaskRequestPayload tests the task request payload handling
func Test_TaskRequestPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name: "valid JSON payload",
			payload: `{
				"type": "weather_check",
				"location": {
					"latitude": 40.7128,
					"longitude": -74.0060,
					"city": "New York",
					"country": "USA"
				},
				"threshold": 25.0
			}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			payload: `{"type": "weather_check", invalid json}`,
			wantErr: true,
		},
		{
			name:    "empty payload",
			payload: `{}`,
			wantErr: true,
		},
		{
			name:    "null payload",
			payload: `null`,
			wantErr: true,
		},
	}

	logger := zap.NewNop()
	tw := NewTaskWorker(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &performerV1.TaskRequest{
				TaskId:  []byte("test-" + tt.name),
				Payload: json.RawMessage(tt.payload),
			}

			err := tw.ValidateTask(task)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}