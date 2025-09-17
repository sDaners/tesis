package integration

import (
	"sql-parser/models"
	"time"
)

// ConversationMessage represents a single message in a conversation
type ConversationMessage struct {
	Role    string `json:"role"`    // "system", "user", or "assistant"
	Content string `json:"content"` // The message content
}

// ConversationSession manages a conversation with OpenAI
type ConversationSession struct {
	ID        string                `json:"id"`
	Messages  []ConversationMessage `json:"messages"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Model     string                `json:"model"`
}

// OpenAIRequest represents a request to the OpenAI API
type OpenAIRequest struct {
	Model       string                `json:"model"`
	Messages    []ConversationMessage `json:"messages"`
	Temperature float64               `json:"temperature,omitempty"`
	MaxTokens   int                   `json:"max_tokens,omitempty"`
}

// OpenAIResponse represents a response from the OpenAI API
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// IterationResult represents the outcome of a single iteration
type IterationResult struct {
	Iteration    int                   `json:"iteration"`
	TestResults  models.TestFileResult `json:"test_results"`
	Success      bool                  `json:"success"`
	GeneratedSQL string                `json:"generated_sql"`
}

// PipelineResult represents the outcome of a complete pipeline run
type PipelineResult struct {
	SessionID        string                `json:"session_id"`
	InitialPrompt    string                `json:"initial_prompt"`
	GeneratedSQL     string                `json:"generated_sql"`
	TestResults      models.TestFileResult `json:"test_results"`
	Iterations       int                   `json:"iterations"`
	IterationResults []IterationResult     `json:"iteration_results"`
	Success          bool                  `json:"success"`
	Messages         []ConversationMessage `json:"messages"`
	TotalTime        time.Duration         `json:"total_time"`
	TokensUsed       int                   `json:"tokens_used"`
}

// OpenAIConfig holds configuration for OpenAI API calls
type OpenAIConfig struct {
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	BaseURL     string  `json:"base_url"`
}
