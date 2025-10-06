package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultOpenAIURL     = "https://api.openai.com/v1/chat/completions"
	ConversationsBaseURL = "https://api.openai.com/v1/conversations"
	DefaultModel         = "chatgpt-4o-latest"
	DefaultTimeout       = 10 * time.Minute
	RetryDelaySeconds    = 30
	MaxRetries           = 10
)

// OpenAIErrorResponse represents an error response from OpenAI API
type OpenAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

func isRateLimitError(body []byte) bool {
	var errorResp OpenAIErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return false
	}
	return errorResp.Error.Code == "rate_limit_exceeded"
}

// isGPT5OrNewer checks if the model uses the new parameter format
func isGPT5OrNewer(model string) bool {
	// Models that use max_completion_tokens instead of max_tokens
	gpt5Models := []string{
		"gpt-5-2025-08-07",
		"gpt-5-mini-2025-08-07",
	}

	for _, gpt5Model := range gpt5Models {
		if model == gpt5Model ||
			strings.HasPrefix(model, gpt5Model+"-") ||
			strings.Contains(model, gpt5Model) {
			return true
		}
	}
	return false
}

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
	if config.MaxTokens == 0 {
		config.MaxTokens = 8096
	}

	return &OpenAIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

func (c *OpenAIClient) SendMessage(messages []ConversationMessage) (*OpenAIResponse, error) {
	request := OpenAIRequest{
		Model:    c.config.Model,
		Messages: messages,
		TopP:     0.5,
	}

	// Use the appropriate token parameter based on the model
	if isGPT5OrNewer(c.config.Model) {
		request.MaxCompletionTokens = c.config.MaxTokens
		request.Temperature = 1
	} else {
		request.MaxTokens = c.config.MaxTokens
		request.Temperature = c.config.Temperature
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Retry logic for rate limit errors
	for attempt := 0; attempt < MaxRetries; attempt++ {
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

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if c.config.Verbose {
			log.Printf("OpenAI API Response Body: %s", string(body))
		}

		if resp.StatusCode == http.StatusOK {
			// Success - parse and return response
			if c.config.Verbose {
				log.Printf("OpenAI API Response Body: %s", string(body))
			}

			var openAIResp OpenAIResponse
			if err := json.Unmarshal(body, &openAIResp); err != nil {
				return nil, fmt.Errorf("failed to unmarshal response: %w", err)
			}

			return &openAIResp, nil
		}

		// Check if this is a rate limit error
		if resp.StatusCode == http.StatusTooManyRequests && isRateLimitError(body) {
			if attempt < MaxRetries-1 {
				log.Printf("Rate limit exceeded, retrying in %d seconds (attempt %d/%d)...",
					RetryDelaySeconds, attempt+1, MaxRetries)
				time.Sleep(RetryDelaySeconds * time.Second)
				continue
			}
		}

		// Not a rate limit error or max retries exceeded - return the error
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil, fmt.Errorf("OpenAI API error: max retries exceeded for rate limit")
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

// Conversations API methods

// CreateConversation creates a new conversation (fallback implementation using chat completions)
func (c *OpenAIClient) CreateConversation() (*CreateConversationResponse, error) {
	// Since the Conversations API might not be available yet, we'll simulate it
	// by generating a conversation ID and managing state locally
	conversationID := fmt.Sprintf("conv_%d", time.Now().Unix())

	return &CreateConversationResponse{
		ID:        conversationID,
		Object:    "conversation",
		CreatedAt: time.Now().Unix(),
	}, nil
}

// AddMessage adds a message to an existing conversation (fallback implementation)
func (c *OpenAIClient) AddMessage(conversationID, role, content string) (*AddMessageResponse, error) {
	// Simulate adding a message by generating a message ID
	messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())

	return &AddMessageResponse{
		ID:             messageID,
		Object:         "message",
		CreatedAt:      time.Now().Unix(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
	}, nil
}

// GetResponse gets an AI response from the conversation using chat completions
func (c *OpenAIClient) GetResponse(conversationID string, messages []ConversationMessage) (*GetResponseResponse, error) {
	response, err := c.SendMessage(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from chat completions: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// Convert to GetResponseResponse format
	responseID := fmt.Sprintf("resp_%d", time.Now().UnixNano())

	return &GetResponseResponse{
		ID:             responseID,
		Object:         "response",
		CreatedAt:      time.Now().Unix(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        response.Choices[0].Message.Content,
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		},
	}, nil
}
