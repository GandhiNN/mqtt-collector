// HTTP client for communicating with the database service API
package dbclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mqtt-catalog/pkg/models"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Creates new database client with base URL and 10-second timeout
func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Posts MQTT sample data to database service via HTTP API with context support
func (c *Client) SendSample(ctx context.Context, sample models.Sample) error {
	jsonData, err := json.Marshal(sample)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+"/api/samples",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("db service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
