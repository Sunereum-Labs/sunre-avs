package main

import (
	"bytes"
	"encoding/json"
	"testing"

	performerV1 "github.com/Layr-Labs/protocol-apis/gen/protos/eigenlayer/hourglass/v1/performer"
	"go.uber.org/zap"
)

func TestSunReWorker_ValidateTask(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	worker := NewSunReWorker(logger)

	tests := []struct {
		name    string
		payload []byte
		wantErr bool
	}{
		{
			name: "valid NYC weather task",
			payload: []byte(`{
				"location": {"latitude": 40.7128, "longitude": -74.0060, "city": "New York"},
				"timestamp": 1704067200,
				"policy_id": "POL-001"
			}`),
			wantErr: false,
		},
		{
			name: "invalid latitude",
			payload: []byte(`{
				"location": {"latitude": 91, "longitude": -74.0060},
				"timestamp": 1704067200,
				"policy_id": "POL-001"
			}`),
			wantErr: true,
		},
		{
			name: "invalid longitude",
			payload: []byte(`{
				"location": {"latitude": 40.7128, "longitude": -181},
				"timestamp": 1704067200,
				"policy_id": "POL-001"
			}`),
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			payload: []byte(`{invalid json}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &performerV1.TaskRequest{
				TaskId:  []byte("test-task-1"),
				Payload: tt.payload,
			}

			err := worker.ValidateTask(task)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSunReWorker_HandleTask(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	worker := NewSunReWorker(logger)

	validPayload := []byte(`{
		"location": {"latitude": 40.7128, "longitude": -74.0060, "city": "New York"},
		"timestamp": 1704067200,
		"policy_id": "POL-001"
	}`)

	task := &performerV1.TaskRequest{
		TaskId:  []byte("test-task-1"),
		Payload: validPayload,
	}

	response, err := worker.HandleTask(task)
	if err != nil {
		t.Fatalf("HandleTask() error = %v", err)
	}

	if !bytes.Equal(response.TaskId, task.TaskId) {
		t.Errorf("Response TaskId = %v, want %v", response.TaskId, task.TaskId)
	}

	// Parse response to verify structure
	var result map[string]interface{}
	if err := json.Unmarshal(response.Result, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check required fields in response
	requiredFields := []string{"task_id", "policy_id", "location", "weather", "verified", "timestamp"}
	for _, field := range requiredFields {
		if _, ok := result[field]; !ok {
			t.Errorf("Response missing required field: %s", field)
		}
	}

	// Verify weather data structure
	if weather, ok := result["weather"].(map[string]interface{}); ok {
		weatherFields := []string{"temperature", "humidity", "wind_speed", "timestamp", "source"}
		for _, field := range weatherFields {
			if _, ok := weather[field]; !ok {
				t.Errorf("Weather data missing field: %s", field)
			}
		}
	} else {
		t.Error("Response missing weather data")
	}
}

func TestWeatherDataGeneration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	worker := NewSunReWorker(logger)

	locations := []Location{
		{Latitude: 40.7128, Longitude: -74.0060, City: "New York"},
		{Latitude: 25.7617, Longitude: -80.1918, City: "Miami"},
		{Latitude: 51.5074, Longitude: -0.1278, City: "London"},
		{Latitude: -33.8688, Longitude: 151.2093, City: "Sydney"},
	}

	for _, loc := range locations {
		weather := worker.generateFallbackWeatherData(loc)
		
		// Verify temperature is reasonable
		if weather.Temperature < -50 || weather.Temperature > 60 {
			t.Errorf("Unrealistic temperature %f for %s", weather.Temperature, loc.City)
		}
		
		// Verify humidity is in valid range
		if weather.Humidity < 0 || weather.Humidity > 100 {
			t.Errorf("Invalid humidity %f for %s", weather.Humidity, loc.City)
		}
		
		// Verify wind speed is reasonable
		if weather.WindSpeed < 0 || weather.WindSpeed > 200 {
			t.Errorf("Invalid wind speed %f for %s", weather.WindSpeed, loc.City)
		}
		
		// Verify timestamp is set
		if weather.Timestamp.IsZero() {
			t.Error("Weather timestamp not set")
		}
		
		// Verify source is set
		if weather.Source == "" {
			t.Error("Weather source not set")
		}
	}
}