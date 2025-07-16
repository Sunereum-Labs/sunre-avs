package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	performerV1 "github.com/Layr-Labs/protocol-apis/gen/protos/eigenlayer/hourglass/v1/performer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// HTTP bridge server to connect UI to gRPC performer
func main() {
	// Connect to gRPC performer
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to performer: %v", err)
	}
	defer conn.Close()

	client := performerV1.NewPerformerServiceClient(conn)

	// CORS middleware
	corsHandler := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			h(w, r)
		}
	}

	// Health endpoint
	http.HandleFunc("/health", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}))

	// Task submission endpoint (for UI)
	http.HandleFunc("/task", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Payload string `json:"payload"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Decode base64 payload
		payloadBytes, err := base64.StdEncoding.DecodeString(req.Payload)
		if err != nil {
			http.Error(w, "Invalid base64 payload", http.StatusBadRequest)
			return
		}

		// Create task request
		taskReq := &performerV1.TaskRequest{
			TaskId:  []byte(fmt.Sprintf("task-%d", time.Now().UnixNano())),
			Payload: payloadBytes,
		}

		// Submit to performer
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := client.ExecuteTask(ctx, taskReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode response
		result := map[string]interface{}{
			"TaskId": string(resp.TaskId),
			"Result": base64.StdEncoding.EncodeToString(resp.Result),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	// DevKit-compatible task submission endpoint
	http.HandleFunc("/submit-task", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Accept raw JSON payload (not base64 encoded)
		var taskData interface{}
		if err := json.NewDecoder(r.Body).Decode(&taskData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Convert to JSON bytes
		payloadBytes, err := json.Marshal(taskData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Received task: %s", string(payloadBytes))

		// Create task request
		taskReq := &performerV1.TaskRequest{
			TaskId:  []byte(fmt.Sprintf("devkit-task-%d", time.Now().UnixNano())),
			Payload: payloadBytes,
		}

		// Submit to performer
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := client.ExecuteTask(ctx, taskReq)
		if err != nil {
			log.Printf("Task submission error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse the result
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Result, &result); err != nil {
			log.Printf("Failed to parse result: %v", err)
			result = map[string]interface{}{
				"raw_result": base64.StdEncoding.EncodeToString(resp.Result),
				"task_id": string(resp.TaskId),
			}
		}

		log.Printf("Task completed successfully: %s", string(resp.Result))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"task_id": string(resp.TaskId),
			"result": result,
		})
	}))

	log.Println("HTTP bridge listening on :8081")
	log.Println("- Task submission (UI): POST /task")
	log.Println("- Task submission (DevKit): POST /submit-task")
	log.Println("- Health check: GET /health")
	log.Fatal(http.ListenAndServe(":8081", nil))
}