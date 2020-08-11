package auth

import (
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Session holds information related to authenticating and
// authorizing the contained user
type Session struct {
	User    *models.User
	Session *models.Session
	Roles   models.UserRoles
}

// NewSession returns a new auth Session
func NewSession(
	user *models.User,
	session *models.Session,
	roles models.UserRoles) (*Session, error) {
	s := &Session{
		User:    user,
		Session: session,
		Roles:   roles,
	}

	return s, s.Validate()
}

// Validate validates that the Session contains valid data
func (s *Session) Validate() error {
	if s.User == nil {
		return errors.New(InvalidInfo, "User must not be nil")
	}

	if s.Session == nil {
		return errors.New(InvalidInfo, "Session must not be nil")
	}

	return nil
}

func (s *Session) GetID() string {
	return s.User.ID
}
