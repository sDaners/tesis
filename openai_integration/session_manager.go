package integration

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// SessionManager manages conversation sessions with OpenAI using the Conversations API
type SessionManager struct {
	sessions map[string]*ConversationSession
	client   *OpenAIClient
	// Keep local message tracking for backward compatibility and reporting
	messages map[string][]ConversationMessage
}

// NewSessionManager creates a new session manager
func NewSessionManager(client *OpenAIClient) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*ConversationSession),
		client:   client,
		messages: make(map[string][]ConversationMessage),
	}
}

// CreateSession creates a new conversation session using the OpenAI Conversations API
func (sm *SessionManager) CreateSession(model string) (*ConversationSession, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	if model == "" {
		model = DefaultModel
	}

	conversationResp, err := sm.client.CreateConversation()
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation on OpenAI: %w", err)
	}

	session := &ConversationSession{
		ID:             sessionID,
		ConversationID: conversationResp.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Model:          model,
		MessageCount:   0,
		LastMessageID:  "",
		LastResponseID: "",
	}

	sm.sessions[sessionID] = session
	sm.messages[sessionID] = make([]ConversationMessage, 0)
	return session, nil
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (*ConversationSession, error) {
	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	return session, nil
}

// AddMessage adds a message to a conversation using the OpenAI Conversations API
func (sm *SessionManager) AddMessage(sessionID string, role string, content string) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}

	// Add the message to the conversation on OpenAI's side
	messageResp, err := sm.client.AddMessage(session.ConversationID, role, content)
	if err != nil {
		return fmt.Errorf("failed to add message to OpenAI conversation: %w", err)
	}

	// Store message locally for backward compatibility
	message := ConversationMessage{
		Role:    role,
		Content: content,
	}
	sm.messages[sessionID] = append(sm.messages[sessionID], message)

	// Update session metadata
	session.MessageCount++
	session.LastMessageID = messageResp.ID
	session.UpdatedAt = time.Now()

	return nil
}

// SendMessage sends a user message and gets AI response using the Conversations API
func (sm *SessionManager) SendMessage(sessionID string, userMessage string) (string, error) {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return "", err
	}

	// Add user message to the conversation
	if err := sm.AddMessage(sessionID, "user", userMessage); err != nil {
		return "", err
	}

	// Get AI response from the conversation by passing the current message history
	currentMessages := sm.messages[sessionID]
	responseResp, err := sm.client.GetResponse(session.ConversationID, currentMessages)
	if err != nil {
		return "", fmt.Errorf("failed to get response from OpenAI conversation: %w", err)
	}

	// Store AI response locally for backward compatibility
	assistantMessage := ConversationMessage{
		Role:    "assistant",
		Content: responseResp.Content,
	}
	sm.messages[sessionID] = append(sm.messages[sessionID], assistantMessage)

	// Update session metadata with response info
	session.LastResponseID = responseResp.ID
	session.MessageCount++
	session.UpdatedAt = time.Now()

	return responseResp.Content, nil
}

// GetConversationHistory returns the conversation history from local storage
func (sm *SessionManager) GetConversationHistory(sessionID string) ([]ConversationMessage, error) {
	_, err := sm.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	messages, exists := sm.messages[sessionID]
	if !exists {
		return []ConversationMessage{}, nil
	}

	return messages, nil
}

// ListSessions returns all active session IDs
func (sm *SessionManager) ListSessions() []string {
	var sessionIDs []string
	for id := range sm.sessions {
		sessionIDs = append(sessionIDs, id)
	}
	return sessionIDs
}

// GetSessionStats returns basic statistics about a session
func (sm *SessionManager) GetSessionStats(sessionID string) (int, int, error) {
	_, err := sm.GetSession(sessionID)
	if err != nil {
		return 0, 0, err
	}

	messages, exists := sm.messages[sessionID]
	if !exists {
		return 0, 0, nil
	}

	userMessages := 0
	assistantMessages := 0

	for _, msg := range messages {
		switch msg.Role {
		case "user":
			userMessages++
		case "assistant":
			assistantMessages++
		}
	}

	return userMessages, assistantMessages, nil
}

// GetConversationID returns the OpenAI conversation ID for a session
func (sm *SessionManager) GetConversationID(sessionID string) (string, error) {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return "", err
	}
	return session.ConversationID, nil
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
