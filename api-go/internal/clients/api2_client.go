// Package clients isolates outbound HTTP calls (API 1 → API 2).
// Keeps services testable via interface mocking.
package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/interseguros/challenge/api-go/internal/config"
	"github.com/interseguros/challenge/api-go/internal/models"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

// StatsClient defines the contract for calling API 2.
type StatsClient interface {
	SendQRForStats(ctx context.Context, payload *models.API2StatsPayload) (map[string]interface{}, error)
}

type httpStatsClient struct {
	baseURL    string
	path       string
	httpClient *http.Client
	logger     utils.Logger
}

// NewStatsClient builds an HTTP client with timeout from config.
func NewStatsClient(cfg *config.Config, logger utils.Logger) StatsClient {
	return &httpStatsClient{
		baseURL: cfg.API2BaseURL,
		path:    cfg.API2MatrixPath,
		httpClient: &http.Client{
			Timeout: cfg.API2Timeout,
		},
		logger: logger,
	}
}

// SendQRForStats POSTs Q and R to API 2 /api/stats.
func (c *httpStatsClient) SendQRForStats(ctx context.Context, payload *models.API2StatsPayload) (map[string]interface{}, error) {
	url := c.baseURL + c.path
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	c.logger.Info("calling api2 stats", "url", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("api2 status %d: %s", resp.StatusCode, string(raw))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
