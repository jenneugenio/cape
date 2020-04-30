package auth

import (
	"github.com/capeprivacy/cape/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Info holds information related to authenticating and
// authorizing the contained identity
type Session struct {
	identity primitives.Identity
	session  *primitives.Session
	policies []*primitives.Policy
}

// NewSession returns a new auth Session
func NewSession(identity primitives.Identity, session *primitives.Session, policies []*primitives.Policy) (*Session, error) {
	info := &Session{
		identity: identity,
		session:  session,
		policies: policies,
	}

	return info, info.Validate()
}

func (s *Session) Validate() error {
	if s.identity == nil {
		return errors.New(InvalidInfo, "Identity must not be nil")
	}

	if s.session == nil {
		return errors.New(InvalidInfo, "Session must not be nil")
	}

	if s.policies == nil {
		return errors.New(InvalidInfo, "Policies must not be nil")
	}

	return nil
}

// Can checks to see if the given identity can do an action on the given primitive type. This
// is intended to work on internal authorization and policy decisions.
func (s *Session) Can(action primitives.Action, typ types.Type) (bool, error) {
	var rules []*primitives.Rule

	for _, p := range s.policies {
		for _, r := range p.Spec.Rules {
			if r.Target.Collection().String() == typ.String() && r.Action == action {
				if r.Effect == primitives.Deny {
					return false, errors.New(AuthorizationFailure, "A rule denies this action")
				}
				rules = append(rules, r)
			}
		}
	}

	if len(rules) == 0 {
		return false, errors.New(AuthorizationFailure, "No rules match this target and action")
	}

	return true, nil
}
