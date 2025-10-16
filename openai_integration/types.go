package integration

import (
	"sql-parser/models"
	"time"
)

type ConversationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ConversationSession struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Model          string    `json:"model"`
	MessageCount   int       `json:"message_count"`
	LastMessageID  string    `json:"last_message_id"`
	LastResponseID string    `json:"last_response_id"`
}

type OpenAIRequest struct {
	Model               string                `json:"model"`
	Messages            []ConversationMessage `json:"messages"`
	Temperature         float64               `json:"temperature,omitempty"`
	TopP                float64               `json:"top_p,omitempty"`
	MaxTokens           int                   `json:"max_tokens,omitempty"`
	MaxCompletionTokens int                   `json:"max_completion_tokens,omitempty"`
}

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

type IterationResult struct {
	Iteration    int                   `json:"iteration"`
	TestResults  models.TestFileResult `json:"test_results"`
	Success      bool                  `json:"success"`
	GeneratedSQL string                `json:"generated_sql"`
}

type PipelineResult struct {
	SessionID        string                `json:"session_id"`
	ConversationID   string                `json:"conversation_id"`
	InitialPrompt    string                `json:"initial_prompt"`
	GeneratedSQL     string                `json:"generated_sql"`
	TestResults      models.TestFileResult `json:"test_results"`
	Iterations       int                   `json:"iterations"`
	IterationResults []IterationResult     `json:"iteration_results"`
	Success          bool                  `json:"success"`
	Messages         []ConversationMessage `json:"messages"`
	TotalTime        time.Duration         `json:"total_time"`
	TokensUsed       int                   `json:"tokens_used"`
	ExecutionMode    string                `json:"execution_mode"`
	Timestamp        time.Time             `json:"timestamp"`
}

type OpenAIConfig struct {
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	BaseURL     string  `json:"base_url"`
	Verbose     bool    `json:"verbose"`
}

type CreateConversationRequest struct {
}

type CreateConversationResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
}

type AddMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AddMessageResponse struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	CreatedAt      int64  `json:"created_at"`
	ConversationID string `json:"conversation_id"`
	Role           string `json:"role"`
	Content        string `json:"content"`
}

type GetResponseRequest struct {
	Model               string  `json:"model,omitempty"`
	Temperature         float64 `json:"temperature,omitempty"`
	MaxTokens           int     `json:"max_tokens,omitempty"`
	MaxCompletionTokens int     `json:"max_completion_tokens,omitempty"`
}

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

type IterationMetrics struct {
	IterationNumber      int     `json:"iteration_number"`
	TotalStatements      int     `json:"total_statements"`
	SuccessfullyParsed   int     `json:"successfully_parsed"`
	ParseErrors          int     `json:"parse_errors"`
	Executed             int     `json:"executed"`
	ExecutionErrors      int     `json:"execution_errors"`
	ParseSuccessRate     float64 `json:"parse_success_rate"`
	ExecutionSuccessRate float64 `json:"execution_success_rate"`
	OverallSuccessRate   float64 `json:"overall_success_rate"`
	Success              bool    `json:"success"`
}

type ExecutionMetrics struct {
	ConversationID     string             `json:"conversation_id"`
	Mode               string             `json:"mode"`
	Model              string             `json:"model"`
	Success            bool               `json:"final_success"`
	TotalIterations    int                `json:"total_iterations"`
	ShortPrompts       bool               `json:"short_prompts"`
	MoreContextEnabled bool               `json:"more_context"`
	IterationResults   []IterationMetrics `json:"iteration_results"`
	Timestamp          time.Time          `json:"timestamp"`
}

type AccumulatedResults struct {
	Executions map[string]*ExecutionMetrics `json:"executions"`
	UpdatedAt  time.Time                    `json:"updated_at"`
	Count      int                          `json:"count"`
}
