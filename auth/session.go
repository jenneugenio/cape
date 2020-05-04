package auth

import (
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Session holds information related to authenticating and
// authorizing the contained identity
type Session struct {
	Identity primitives.Identity
	Session  *primitives.Session
	Policies []*primitives.Policy
	Roles    []*primitives.Role
}

// NewSession returns a new auth Session
func NewSession(identity primitives.Identity, session *primitives.Session, policies []*primitives.Policy,
	roles []*primitives.Role) (*Session, error) {
	s := &Session{
		Identity: identity,
		Session:  session,
		Policies: policies,
		Roles:    roles,
	}

	return s, s.Validate()
}

// Validate validates that the Session contains valid data
func (s *Session) Validate() error {
	if s.Identity == nil {
		return errors.New(InvalidInfo, "Identity must not be nil")
	}

	if s.Session == nil {
		return errors.New(InvalidInfo, "Session must not be nil")
	}

	if s.Policies == nil {
		return errors.New(InvalidInfo, "Policies must not be nil")
	}

	if s.Roles == nil {
		return errors.New(InvalidInfo, "Roles must not be nil")
	}

	return nil
}

// Can checks to see if the given identity can do an action on the given primitive type. This
// is intended to work on internal authorization and policy decisions.
func (s *Session) Can(action primitives.Action, typ types.Type) error {
	var rules []*primitives.Rule

	for _, p := range s.Policies {
		for _, r := range p.Spec.Rules {
			if r.Target.Type().String() == typ.String() && r.Action == action {
				if r.Effect == primitives.Deny {
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
