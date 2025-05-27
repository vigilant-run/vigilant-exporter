package sender

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
	ErrSendInvalidRequest = errors.New("invalid request")
	ErrSendFailed         = errors.New("failed to send batch")
	ErrSendTimeout        = errors.New("send timeout")
	ErrSendCanceled       = errors.New("send canceled")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPSender is a sender that sends logs to a HTTP endpoint.
type HTTPSender struct {
	httpClient HTTPClient

	endpoint string
	token    string

	mux sync.Mutex
}

// NewHTTPSender creates a new HTTPSender.
func NewHTTPSender(
	httpClient HTTPClient,
	endpoint string,
	token string,
) *HTTPSender {
	return &HTTPSender{
		httpClient: httpClient,
		endpoint:   endpoint,
		token:      token,
		mux:        sync.Mutex{},
	}
}

// SendBatch synchronously sends a batch of logs to the HTTP endpoint.
func (e *HTTPSender) SendBatch(
	ctx context.Context,
	batch *data.MessageBatch,
) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	json, err := json.Marshal(batch)
	if err != nil {
		return ErrSendInvalidRequest
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		e.endpoint,
		bytes.NewBuffer(json),
	)
	if err != nil {
		return ErrSendInvalidRequest
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return ErrSendFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrSendFailed
	}

	return nil
}
