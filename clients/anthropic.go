package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
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

func (c *AnthropicClient) Analyze(prompt string, context map[string]interface{}, outputJson string) (interface{}, error) {
	endpoint := "https://api.anthropic.com/v1/messages"

	ctxJson, _ := json.Marshal(context)
	messages := []Message{
		{Role: "user", Content: prompt},
		{Role: "user", Content: "outputJson: " + outputJson},
		{Role: "user", Content: "Context: " + string(ctxJson)},
	}

	reqBody := AnalysisRequest{
		Model:     "claude-3-7-sonnet-20250219",
		Messages:  messages,
		MaxTokens: 3000,
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

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, response)
	}

	content := response["content"].([]interface{})
	output := content[0].(map[string]interface{})
	outputText := output["text"].(string)
	outputText = strings.Replace(outputText, "```json", "", -1)
	outputText = strings.Replace(outputText, "```", "", -1)

	outputMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(outputText), &outputMap); err != nil {
		return nil, fmt.Errorf("error unmarshaling output text: %v", err)
	}
	return outputMap, nil
}

func (c *AnthropicClient) ConvertResponse(input string, outputType reflect.Type) (interface{}, error) {
	outputJson, err := json.MarshalIndent(reflect.New(outputType).Interface(), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling output type: %v", err)
	}

	prompt := fmt.Sprintf("You are a data conversion assistant. Convert the input data to the specified JSON format. Output only the JSON. \n Convert this data into the following JSON format: %s\n\nInput data:\n%s", string(outputJson), input)

	messages := []Message{
		{Role: "user", Content: prompt},
	}

	reqBody := CompletionRequest{
		Model:     "claude-3-7-sonnet-20250219",
		Messages:  messages,
		MaxTokens: 5000,
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
