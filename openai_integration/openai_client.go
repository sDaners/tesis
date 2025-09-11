package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sql-parser/tools"
	"time"
)

const (
	DefaultOpenAIURL = "https://api.openai.com/v1/chat/completions"
	DefaultModel     = "gpt-4o-mini"
	DefaultTimeout   = 60 * time.Second
)

// OpenAIClient handles communication with OpenAI API
type OpenAIClient struct {
	config     OpenAIConfig
	httpClient *http.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config OpenAIConfig) *OpenAIClient {
	if config.BaseURL == "" {
		config.BaseURL = DefaultOpenAIURL
	}
	if config.Model == "" {
		config.Model = DefaultModel
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4096
	}

	return &OpenAIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewOpenAIClientFromEnv creates a new OpenAI client using environment variables
func NewOpenAIClientFromEnv() (*OpenAIClient, error) {
	apiKey := tools.Get().OpenAIAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	config := OpenAIConfig{
		APIKey:      apiKey,
		Model:       DefaultModel,
		Temperature: 0.7,
		MaxTokens:   4096,
		BaseURL:     DefaultOpenAIURL,
	}

	return NewOpenAIClient(config), nil
}

// SendMessage sends a message to OpenAI and returns the response
func (c *OpenAIClient) SendMessage(messages []ConversationMessage) (*OpenAIResponse, error) {
	request := OpenAIRequest{
		Model:       c.config.Model,
		Messages:    messages,
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &openAIResp, nil
}

// SendSingleMessage is a convenience method for sending a single user message
func (c *OpenAIClient) SendSingleMessage(content string) (string, error) {
	messages := []ConversationMessage{
		{
			Role:    "user",
			Content: content,
		},
	}

	response, err := c.SendMessage(messages)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}

// GetUsageFromResponse extracts token usage from OpenAI response
func (c *OpenAIClient) GetUsageFromResponse(response *OpenAIResponse) (int, int, int) {
	return response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens
}
