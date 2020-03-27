package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Service represents a service connecting to the system (e.g. a machine
// running a pipeline).
type Service struct {
	*IdentityImpl
	Type ServiceType `json:"type"`
}

// GetType returns the type for this entity
func (s *Service) GetType() types.Type {
	return ServicePrimitiveType
}

// NewService returns a mutable service struct
func NewService(email Email, typ ServiceType, creds *Credentials) (*Service, error) {
	p, err := database.NewPrimitive(ServicePrimitiveType)
	if err != nil {
		return nil, err
	}

	return &Service{
		IdentityImpl: &IdentityImpl{
			Primitive:   p,
			Email:       email,
			Credentials: creds,
		},
		Type: typ,
	}, nil
}

// GetCredentials satisfies Identity interface
func (s *Service) GetCredentials() *Credentials {
	return s.Credentials
}

// GetEmail satisfies Identity interface
func (s *Service) GetEmail() Email {
	return s.Email
}
