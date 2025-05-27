package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Batch struct {
	Token string `json:"token"`
	Logs  []Log  `json:"logs"`
}

type Log struct {
	Timestamp  time.Time         `json:"timestamp"`
	Level      string            `json:"level"`
	Body       string            `json:"body"`
	Attributes map[string]string `json:"attributes"`
}

type MockServer struct {
	receivedLogs []Log
	mu           sync.RWMutex
}

func (s *MockServer) handleIngress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var batch Batch
	if err := json.Unmarshal(body, &batch); err != nil {
		http.Error(w, "failed to unmarshal body", http.StatusBadRequest)
		return
	}

	if batch.Token != "integration_token" {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	s.mu.Lock()
	s.receivedLogs = append(s.receivedLogs, batch.Logs...)
	s.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *MockServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *MockServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	logs := make([]Log, len(s.receivedLogs))
	copy(logs, s.receivedLogs)
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"count": len(logs),
		"logs":  logs,
	})
}

func main() {
	server := &MockServer{
		receivedLogs: make([]Log, 0),
	}

	http.HandleFunc("/api/message", server.handleIngress)
	http.HandleFunc("/api/health", server.handleHealth)
	http.HandleFunc("/api/logs", server.handleLogs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
