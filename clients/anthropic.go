package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type AnthropicClient struct {
	APIKey string
	Model  string
}

type CompletionRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		APIKey: apiKey,
	}
}

type AnalysisRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
}

type AnalysisResponse struct {
	Content map[string]string `json:"content"`
	Error   string            `json:"error,omitempty"`
}

func (c *AnthropicClient) Analyze(prompt string, context []string, outputJson string) (map[string]string, error) {
	endpoint := "https://api.anthropic.com/v1/messages"

	systemMsg := fmt.Sprintf("You are an analysis assistant. Analyze the data and respond in this JSON format: %s", outputJson)

	messages := []Message{
		{Role: "system", Content: systemMsg},
		{Role: "user", Content: prompt},
	}

	for _, ctx := range context {
		messages = append(messages, Message{Role: "user", Content: ctx})
	}

	reqBody := AnalysisRequest{
		Model:     c.Model,
		Messages:  messages,
		MaxTokens: 1000,
		System:    systemMsg,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var response AnalysisResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, response.Error)
	}

	return response.Content, nil
}

func (c *AnthropicClient) ConvertResponse(input string, outputType reflect.Type) (interface{}, error) {
	outputJson, err := json.MarshalIndent(reflect.New(outputType).Interface(), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling output type: %v", err)
	}

	prompt := fmt.Sprintf("Convert this data into the following JSON format: %s\n\nInput data:\n%s", string(outputJson), input)

	systemMsg := "You are a data conversion assistant. Convert the input data to the specified JSON format. Output only the JSON."

	messages := []Message{
		{Role: "system", Content: systemMsg},
		{Role: "user", Content: prompt},
	}

	reqBody := CompletionRequest{
		Model:     c.Model,
		Messages:  messages,
		MaxTokens: 1000,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	output := reflect.New(outputType).Interface()
	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return output, nil
}