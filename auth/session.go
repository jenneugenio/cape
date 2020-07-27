package auth

import (
	"github.com/capeprivacy/cape/coordinator/database/types"
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Session holds information related to authenticating and
// authorizing the contained user
type Session struct {
	User               *models.User
	Session            *primitives.Session
	Policies           []*models.RBACPolicy
	Roles              models.UserRoles
	CredentialProvider primitives.CredentialProvider
}

// NewSession returns a new auth Session
func NewSession(
	user *models.User,
	session *primitives.Session,
	policies []*models.RBACPolicy,
	roles models.UserRoles,
	cp primitives.CredentialProvider) (*Session, error) {
	s := &Session{
		User:               user,
		Session:            session,
		Policies:           policies,
		Roles:              roles,
		CredentialProvider: cp,
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

	if s.Policies == nil {
		return errors.New(InvalidInfo, "Policies must not be nil")
	}

	return nil
}

func (s *Session) GetID() string {
	return s.User.ID
}

// Can checks to see if the given user can do an action on the given primitive type. This
// is intended to work on internal authorization and policy decisions.
func (s *Session) Can(action models.RBACAction, typ types.Type) error {
	var rules []*models.RBACRule

	for _, p := range s.Policies {
		for _, r := range p.Spec.Rules {
			if r.Target.Type().String() == typ.String() && r.Action == action {
				if r.Effect == models.Deny {
					return errors.New(AuthorizationFailure, "A rule denies this action")
				}
				rules = append(rules, r)
			}
		}
	}

	if len(rules) == 0 {
		return errors.New(AuthorizationFailure, "You don't have sufficient permissions to perform a %s on %s", action, typ)
	}

	return nil
}
