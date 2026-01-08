package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// BaseURL is the seats.aero partner API base URL
	BaseURL = "https://seats.aero/partnerapi"

	// DefaultTimeout for HTTP requests
	DefaultTimeout = 30 * time.Second
)

// Client is the seats.aero API client
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiKey:  apiKey,
		baseURL: BaseURL,
	}
}

// WithTimeout sets a custom timeout for the client
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

// WithBaseURL sets a custom base URL (useful for testing)
func (c *Client) WithBaseURL(url string) *Client {
	c.baseURL = url
	return c
}

// doRequest performs an authenticated HTTP request
func (c *Client) doRequest(method, endpoint string, params map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add auth header
	req.Header.Set("Partner-Authorization", c.apiKey)
	req.Header.Set("Accept", "application/json")

	// Add query parameters
	if len(params) > 0 {
		q := req.URL.Query()
		for key, value := range params {
			if value != "" {
				q.Add(key, value)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// get performs a GET request and unmarshals the response
func (c *Client) get(endpoint string, params map[string]string, result interface{}) error {
	body, err := c.doRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

// APIError represents an error from the seats.aero API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}
