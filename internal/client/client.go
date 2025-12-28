package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL        = "https://console.sannti.cloud/restapi"
	DefaultTimeout = 30 * time.Second
)

// Client represents the Sannti API client
type Client struct {
	HTTPClient *http.Client
	APIKey     string
	SecretKey  string
	BaseURL    string
}

// NewClient creates a new Sannti API client
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   BaseURL,
	}
}

// DoRequest performs an HTTP request with authentication headers
func (c *Client) DoRequest(method, path string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication headers
	req.Header.Set("apikey", c.APIKey)
	req.Header.Set("secretkey", c.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for API errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get performs a GET request
func (c *Client) Get(path string) ([]byte, error) {
	return c.DoRequest("GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}) ([]byte, error) {
	return c.DoRequest("POST", path, body)
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}) ([]byte, error) {
	return c.DoRequest("PUT", path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) ([]byte, error) {
	return c.DoRequest("DELETE", path, nil)
}
