package export

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
	"vigilant-exporter/internal/data"
)

type MockHTTPClient struct {
	DoFunc    func(req *http.Request) (*http.Response, error)
	CallCount int
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.CallCount++
	return m.DoFunc(req)
}

func TestNewHTTPExporter(t *testing.T) {
	client := &http.Client{}
	endpoint := "https://api.example.com/logs"
	token := "test-token"

	exporter := NewHTTPExporter(client, endpoint, token)

	if exporter == nil {
		t.Fatal("Expected exporter to be non-nil")
	}

	if exporter.httpClient != client {
		t.Error("Expected httpClient to be set correctly")
	}

	if exporter.endpoint != endpoint {
		t.Errorf("Expected endpoint to be %s, got %s", endpoint, exporter.endpoint)
	}

	if exporter.token != token {
		t.Errorf("Expected token to be %s, got %s", token, exporter.token)
	}
}

func TestHTTPExporter_ExportBatch_Success(t *testing.T) {
	client := createValidMockClient(t)

	exporter := NewHTTPExporter(
		client,
		"https://api.example.com/logs",
		"test-token",
	)

	batch := data.NewMessageBatch(
		"test-token",
		[]*data.Log{
			data.NewLog(
				time.Now(),
				data.LogLevelInfoName,
				"Test log message 1",
				map[string]string{"key1": "value1"},
			),
		},
	)

	err := exporter.ExportBatch(context.Background(), batch)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if client.CallCount != 1 {
		t.Errorf("Expected 1 call to client, got %d", client.CallCount)
	}
}

func TestHTTPExporter_ExportBatch_Empty(t *testing.T) {
	client := createValidMockClient(t)

	exporter := NewHTTPExporter(
		client,
		"https://api.example.com/logs",
		"test-token",
	)

	batch := data.NewMessageBatch("test-token", []*data.Log{})
	err := exporter.ExportBatch(context.Background(), batch)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if client.CallCount != 1 {
		t.Errorf("Expected 1 call to client, got %d", client.CallCount)
	}
}

func TestHTTPExporter_ExportBatch_Concurrent(t *testing.T) {
	client := createValidMockClient(t)

	exporter := NewHTTPExporter(
		client,
		"https://api.example.com/logs",
		"test-token",
	)

	wg := sync.WaitGroup{}
	for range 10 {
		wg.Add(1)
		go func() {
			batch := data.NewMessageBatch("test-token", []*data.Log{})
			exporter.ExportBatch(context.Background(), batch)
			wg.Done()
		}()
	}

	wg.Wait()
	if client.CallCount != 10 {
		t.Errorf("Expected 1 call to client, got %d", client.CallCount)
	}
}

func TestHTTPExporter_ExportBatch_Unavailable(t *testing.T) {
	client := createUnavailableMockClient(t)

	exporter := NewHTTPExporter(
		client,
		"https://api.example.com/logs",
		"test-token",
	)

	batch := data.NewMessageBatch("test-token", []*data.Log{})
	err := exporter.ExportBatch(context.Background(), batch)
	if err != ErrExportFailed {
		t.Errorf("Expected ErrExportFailed, got %v", err)
	}

	if client.CallCount != 1 {
		t.Errorf("Expected 1 call to client, got %d", client.CallCount)
	}
}

func createValidMockClient(t *testing.T) *MockHTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("Expected POST method, got %s", req.Method)
			}

			contentType := req.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Errorf("Failed to read request body: %v", err)
			}

			var batch data.MessageBatch
			if err := json.Unmarshal(body, &batch); err != nil {
				t.Errorf("Request body is not valid JSON: %v", err)
			}

			return createSuccessResponse(), nil
		},
	}
}

func createSuccessResponse() *http.Response {
	return &http.Response{StatusCode: http.StatusOK,
		Body:   io.NopCloser(strings.NewReader("{}")),
		Header: make(http.Header),
	}
}

func createUnavailableMockClient(t *testing.T) *MockHTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader("{}")),
				Header:     make(http.Header),
			}, nil
		},
	}
}
