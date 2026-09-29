package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	return c.doRaw(method, path, "application/json", bodyReader)
}

func (c *Client) doRaw(method, path, contentType string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
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
		var apiErr struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.Error != "" {
			return nil, fmt.Errorf("api %d: %s", resp.StatusCode, apiErr.Error)
		}
		return nil, fmt.Errorf("api %d: %s", resp.StatusCode, string(respBody))
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

// PostBytes sends a POST request with an arbitrary content type.
func (c *Client) PostBytes(path, contentType string, body []byte) ([]byte, error) {
	return c.doRaw(http.MethodPost, path, contentType, bytes.NewReader(body))
}

// PostMultipart sends a multipart/form-data POST with fields and one optional file.
func (c *Client) PostMultipart(path string, fields map[string]string, fileField, filename string, file []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("write multipart field %s: %w", k, err)
		}
	}
	if fileField != "" {
		fw, err := w.CreateFormFile(fileField, filename)
		if err != nil {
			return nil, fmt.Errorf("create multipart file: %w", err)
		}
		if _, err := fw.Write(file); err != nil {
			return nil, fmt.Errorf("write multipart file: %w", err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}
	return c.doRaw(http.MethodPost, path, w.FormDataContentType(), &buf)
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
