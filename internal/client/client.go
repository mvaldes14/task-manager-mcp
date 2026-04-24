package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client for the doit task manager API.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// New returns a Client configured with the given base URL and API key.
func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body any) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-2xx response %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get sends a GET request to the given path and returns the response body.
func (c *Client) Get(path string) ([]byte, error) {
	return c.do(http.MethodGet, path, nil)
}

// Post sends a POST request to the given path with the given body and returns the response body.
func (c *Client) Post(path string, body any) ([]byte, error) {
	return c.do(http.MethodPost, path, body)
}

// Patch sends a PATCH request to the given path with the given body and returns the response body.
func (c *Client) Patch(path string, body any) ([]byte, error) {
	return c.do(http.MethodPatch, path, body)
}

// Put sends a PUT request to the given path with the given body and returns the response body.
func (c *Client) Put(path string, body any) ([]byte, error) {
	return c.do(http.MethodPut, path, body)
}

// Delete sends a DELETE request to the given path and returns the response body.
func (c *Client) Delete(path string) ([]byte, error) {
	return c.do(http.MethodDelete, path, nil)
}
