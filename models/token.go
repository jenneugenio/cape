package models

import (
	"fmt"
	"github.com/manifoldco/go-base64"
)

type Token struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	// We never want to send Credentials over the wire!
	Credentials *Credentials `json:"-" gqlgen:"-"`
}

type EncryptedToken struct {
	*Token
	Credentials *base64.Value `json:"credentials"`
}

func (tc *Token) Validate() error {
	if tc.UserID == "" {
		return fmt.Errorf("user id must be non-nil")
	}

	if tc.Credentials == nil {
		return fmt.Errorf("credentials must be non-nil")
	}

	return nil
}

func (tc *Token) GetUserID() string {
	return tc.UserID
}

func (tc *Token) GetCredentials() (*Credentials, error) {
	return tc.Credentials, nil
}

func NewToken(userID string, creds *Credentials) Token {
	return Token{
		ID:          NewID(),
		UserID:      userID,
		Credentials: creds,
	}
}

func (tc *Token) GetStringID() string {
	return tc.ID
}

// GenerateToken returns an instantiated token for use in unit testing.
//
// This function _should only ever_ be used inside of a test.
func GenerateToken(user User) (Password, Token) {
	password := GeneratePassword()
	c := GenerateCredentials()

	token := NewToken(user.ID, c)
	return password, token
}
