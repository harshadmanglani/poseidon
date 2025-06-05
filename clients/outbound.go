package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OutboundClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string) *OutboundClient {
	return &OutboundClient{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (c *OutboundClient) Get(path string, response interface{}) error {
	return c.doRequest(http.MethodGet, path, nil, response)
}

func (c *OutboundClient) Post(path string, request interface{}, response interface{}) error {
	return c.doRequest(http.MethodPost, path, request, response)
}

func (c *OutboundClient) doRequest(method, path string, request interface{}, response interface{}) error {
	var body io.Reader
	if request != nil {
		jsonData, err := json.Marshal(request)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}