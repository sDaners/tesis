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

// ConversationSession manages a conversation with OpenAI using the Conversations API
type ConversationSession struct {
	ID             string    `json:"id"`              // Local session ID for tracking
	ConversationID string    `json:"conversation_id"` // OpenAI conversation ID
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Model          string    `json:"model"`
	MessageCount   int       `json:"message_count"`    // Number of messages exchanged
	LastMessageID  string    `json:"last_message_id"`  // ID of the last message
	LastResponseID string    `json:"last_response_id"` // ID of the last AI response
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
	ExecutionMode    string                `json:"execution_mode"` // "single" or "iterative"
	Timestamp        time.Time             `json:"timestamp"`      // When this execution occurred
}

// OpenAIConfig holds configuration for OpenAI API calls
type OpenAIConfig struct {
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	BaseURL     string  `json:"base_url"`
}

// Conversations API types

// CreateConversationRequest represents a request to create a new conversation
type CreateConversationRequest struct {
	// Empty struct - the API might not require parameters for conversation creation
}

// CreateConversationResponse represents the response from creating a conversation
type CreateConversationResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
}

// AddMessageRequest represents a request to add a message to a conversation
type AddMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AddMessageResponse represents the response from adding a message
type AddMessageResponse struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	CreatedAt      int64  `json:"created_at"`
	ConversationID string `json:"conversation_id"`
	Role           string `json:"role"`
	Content        string `json:"content"`
}

// GetResponseRequest represents a request to get a response from the conversation
type GetResponseRequest struct {
	Model       string  `json:"model,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

// GetResponseResponse represents the AI's response in a conversation
type GetResponseResponse struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	CreatedAt      int64  `json:"created_at"`
	ConversationID string `json:"conversation_id"`
	Role           string `json:"role"`
	Content        string `json:"content"`
	Usage          struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// IterationMetrics contains metrics for a single iteration
type IterationMetrics struct {
	IterationNumber      int     `json:"iteration_number"`
	TotalStatements      int     `json:"total_statements"`
	SuccessfullyParsed   int     `json:"successfully_parsed"`
	ParseErrors          int     `json:"parse_errors"`
	Executed             int     `json:"executed"`
	ExecutionErrors      int     `json:"execution_errors"`
	ParseSuccessRate     float64 `json:"parse_success_rate"`     // %
	ExecutionSuccessRate float64 `json:"execution_success_rate"` // % of parsed
	OverallSuccessRate   float64 `json:"overall_success_rate"`   // %
	Success              bool    `json:"success"`                // true if no errors
}

// ExecutionMetrics contains metrics for an entire execution with all iterations
type ExecutionMetrics struct {
	ConversationID   string             `json:"conversation_id"`
	Mode             string             `json:"mode"`          // "single" or "iterative"
	Success          bool               `json:"final_success"` // final outcome
	TotalIterations  int                `json:"total_iterations"`
	IterationResults []IterationMetrics `json:"iteration_results"` // metrics for each iteration
	Timestamp        time.Time          `json:"timestamp"`
}

// AccumulatedResults stores all pipeline executions for analysis and graphing
type AccumulatedResults struct {
	Executions map[string]*ExecutionMetrics `json:"executions"` // Key: conversation_id
	UpdatedAt  time.Time                    `json:"updated_at"`
	Count      int                          `json:"count"` // Total number of executions stored
}
