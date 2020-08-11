package models

import (
	"fmt"
	"time"

	"github.com/manifoldco/go-base64"
)

// Session holds all the session data required to authenticate API
// calls with the server
type Session struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	OwnerID   string        `json:"owner_id"`
	ExpiresAt time.Time     `json:"expires_at"`
	Token     *base64.Value `json:"token"`
}

func (s *Session) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("id is empty")
	}

	if s.UserID == "" {
		return fmt.Errorf("user id is empty")
	}

	if s.OwnerID == "" {
		return fmt.Errorf("owner id is empty")
	}

	return nil
}

// NewSession returns a new Session struct
func NewSession(cp CredentialProvider) Session {
	return Session{
		ID:      NewID(),
		UserID:  cp.GetUserID(),
		OwnerID: cp.GetStringID(),
	}
}

func (s *Session) SetToken(token *base64.Value, expiresAt time.Time) {
	s.Token = token
	s.ExpiresAt = expiresAt
}
