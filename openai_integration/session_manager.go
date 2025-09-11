package integration

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// SessionManager manages conversation sessions with OpenAI
type SessionManager struct {
	sessions map[string]*ConversationSession
	client   *OpenAIClient
}

// NewSessionManager creates a new session manager
func NewSessionManager(client *OpenAIClient) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*ConversationSession),
		client:   client,
	}
}

// CreateSession creates a new conversation session
func (sm *SessionManager) CreateSession(model string) (*ConversationSession, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	if model == "" {
		model = DefaultModel
	}

	session := &ConversationSession{
		ID:        sessionID,
		Messages:  make([]ConversationMessage, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Model:     model,
	}

	sm.sessions[sessionID] = session
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

// AddMessage adds a message to a session
func (sm *SessionManager) AddMessage(sessionID string, role string, content string) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}

	message := ConversationMessage{
		Role:    role,
		Content: content,
	}

	session.Messages = append(session.Messages, message)
	session.UpdatedAt = time.Now()

	return nil
}

// SendMessage sends the current conversation to OpenAI and adds the response
func (sm *SessionManager) SendMessage(sessionID string, userMessage string) (string, error) {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return "", err
	}

	// Add user message to session
	if err := sm.AddMessage(sessionID, "user", userMessage); err != nil {
		return "", err
	}

	// Send all messages in the session to OpenAI
	response, err := sm.client.SendMessage(session.Messages)
	if err != nil {
		return "", fmt.Errorf("failed to send message to OpenAI: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	assistantResponse := response.Choices[0].Message.Content

	// Add assistant response to session
	if err := sm.AddMessage(sessionID, "assistant", assistantResponse); err != nil {
		return "", err
	}

	return assistantResponse, nil
}

// AddSystemMessage adds a system message to the session (should be called first)
func (sm *SessionManager) AddSystemMessage(sessionID string, content string) error {
	return sm.AddMessage(sessionID, "system", content)
}

// GetConversationHistory returns the full conversation history for a session
func (sm *SessionManager) GetConversationHistory(sessionID string) ([]ConversationMessage, error) {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	return session.Messages, nil
}

// DeleteSession removes a session
func (sm *SessionManager) DeleteSession(sessionID string) error {
	if _, exists := sm.sessions[sessionID]; !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	delete(sm.sessions, sessionID)
	return nil
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
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return 0, 0, err
	}

	userMessages := 0
	assistantMessages := 0

	for _, msg := range session.Messages {
		switch msg.Role {
		case "user":
			userMessages++
		case "assistant":
			assistantMessages++
		}
	}

	return userMessages, assistantMessages, nil
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
