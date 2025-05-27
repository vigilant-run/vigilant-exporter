package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const (
	EXECUTABLE_PATH = "/app/vigilant-exporter"

	ENDPOINT_HEALTH  = "http://server:8000/api/health"
	ENDPOINT_LOGS    = "http://server:8000/api/logs"
	ENDPOINT_MESSAGE = "http://server:8000/api/message"

	LOG_DIR       = "/logs"
	LOG_FILE      = "test.log"
	LOG_FILE_PATH = LOG_DIR + "/" + LOG_FILE

	TOKEN    = "integration_token"
	INSECURE = true
)

var testLogs = []string{
	"2024-01-01T10:00:00Z INFO Application started",
	"2024-01-01T10:00:01Z DEBUG Initializing components",
	"2024-01-01T10:00:02Z INFO Server listening on port 8080",
}

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

func TestCLI_ExistingFile(t *testing.T) {
	if err := getHealthFromServer(t); err != nil {
		t.Fatalf("failed to get health: %v", err)
	}

	file, err := createLogFile(t, LOG_DIR, LOG_FILE)
	if err != nil {
		t.Fatalf("failed to create test log file: %v", err)
	}

	if err := writeLogsToFile(file, testLogs); err != nil {
		t.Fatalf("failed to write logs to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("failed to close test log file: %v", err)
	}

	cmd := runCli(t, LOG_FILE_PATH)
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	time.Sleep(2 * time.Second)

	if err := checkLogsCount(t, len(testLogs)); err != nil {
		t.Fatalf("failed to check logs count: %v", err)
	}

	if err := checkLogsReceived(t, testLogs); err != nil {
		t.Fatalf("failed to check logs received: %v", err)
	}
}

func createLogFile(t *testing.T, dir string, filename string) (*os.File, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create log directory: %v", err)
	}

	logFile := filepath.Join(dir, filename)
	file, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("failed to create test log file: %v", err)
	}

	return file, nil
}

func writeLogsToFile(file *os.File, logs []string) error {
	for _, logLine := range logs {
		if _, err := file.WriteString(logLine + "\n"); err != nil {
			return fmt.Errorf("failed to write to log file: %v", err)
		}
	}

	return nil
}

func runCli(t *testing.T, logFile string) *exec.Cmd {
	cmd := exec.Command(EXECUTABLE_PATH,
		"--file", logFile,
		"--token", TOKEN,
		"--endpoint", ENDPOINT_MESSAGE,
		"--insecure",
	)

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start cli: %v", err)
	}

	return cmd
}

func checkLogsCount(t *testing.T, expectedCount int) error {
	logResponse, err := getLogsFromServer(t)
	if err != nil {
		t.Fatalf("failed to get logs from server: %v", err)
	}

	if logResponse.Count != expectedCount {
		return fmt.Errorf("expected %d logs, got %d", expectedCount, logResponse.Count)
	}

	return nil
}

func checkLogsReceived(t *testing.T, expectedLogs []string) error {
	logResponse, err := getLogsFromServer(t)
	if err != nil {
		t.Fatalf("failed to get logs from server: %v", err)
	}

	for _, expectedLog := range expectedLogs {
		found := false
		for _, log := range logResponse.Logs {
			if log.Body == expectedLog {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("expected log not found in received logs: %s", expectedLog)
		}
	}

	return nil
}

func getLogsFromServer(t *testing.T) (*LogResponse, error) {
	resp, err := http.Get(ENDPOINT_LOGS)
	if err != nil {
		t.Fatalf("failed to get logs from server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to get logs: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var logResponse LogResponse
	if err := json.Unmarshal(body, &logResponse); err != nil {
		t.Fatalf("failed to unmarshal log response: %v", err)
	}

	return &logResponse, nil
}

func getHealthFromServer(t *testing.T) error {
	resp, err := http.Get(ENDPOINT_HEALTH)
	if err != nil {
		t.Fatalf("failed to get health from server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server not healthy: %v", resp.Status)
	}

	return nil
}
