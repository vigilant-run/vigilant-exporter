package export

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"vigilant-exporter/internal/data"
)

var (
	ErrExportInvalidRequest = errors.New("invalid request")
	ErrExportFailed         = errors.New("failed to export batch")
	ErrExportTimeout        = errors.New("export timeout")
	ErrExportCanceled       = errors.New("export canceled")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPExporter is an exporter that sends logs to a HTTP endpoint.
type HTTPExporter struct {
	httpClient HTTPClient

	endpoint string
	token    string

	mux sync.Mutex
}

// NewHTTPExporter creates a new HTTPExporter.
func NewHTTPExporter(
	httpClient HTTPClient,
	endpoint string,
	token string,
) *HTTPExporter {
	return &HTTPExporter{
		httpClient: httpClient,
		endpoint:   endpoint,
		token:      token,
		mux:        sync.Mutex{},
	}
}

// ExportBatch synchronously exports a batch of logs to the HTTP endpoint.
func (e *HTTPExporter) ExportBatch(
	ctx context.Context,
	batch *data.MessageBatch,
) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	json, err := json.Marshal(batch)
	if err != nil {
		return ErrExportInvalidRequest
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		e.endpoint,
		bytes.NewBuffer(json),
	)
	if err != nil {
		return ErrExportInvalidRequest
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return ErrExportFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrExportFailed
	}

	return nil
}
