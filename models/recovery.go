package models

import (
	"fmt"
	"time"
)

// RecoveryExpiration is the amount of time that has passed since a recovery
// was created before it's no longer valid.
var RecoveryExpiration = 30 * time.Minute

type Recovery struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	Credentials *Credentials `json:"-" gqlgen:"-"`
	ExpiresAt   time.Time    `json:"expires_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (r *Recovery) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("id must not be empty")
	}

	if r.UserID == "" {
		return fmt.Errorf("user id must not be empty")
	}

	if r.Credentials == nil {
		return fmt.Errorf("missing credentials")
	}

	if r.ExpiresAt.IsZero() {
		return fmt.Errorf("missing expires at")
	}

	return nil
}

func (r *Recovery) Expired() bool {
	return time.Now().UTC().After(r.ExpiresAt)
}

func NewRecovery(userID string, creds *Credentials) Recovery {
	r := Recovery{
		ID:          NewID(),
		UserID:      userID,
		Credentials: creds,
		ExpiresAt:   time.Now().UTC().Add(RecoveryExpiration),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return r
}

func GenerateRecovery() Recovery {
	userID := "thisisanid"
	return NewRecovery(userID, GenerateCredentials())
}
