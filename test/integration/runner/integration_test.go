package tests

import (
	"net/http"
	"testing"
	"time"
)

type Log struct {
	Timestamp  time.Time         `json:"timestamp"`
	Level      string            `json:"level"`
	Body       string            `json:"body"`
	Attributes map[string]string `json:"attributes"`
}

type LogResponse struct {
	Count int   `json:"count"`
	Logs  []Log `json:"logs"`
}

func TestRealTimeLogTailing(t *testing.T) {
	resp, err := http.Get("http://server:8000/api/health")
	if err != nil {
		t.Fatalf("failed to get health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to get health: %v", resp.Status)
	}
}
