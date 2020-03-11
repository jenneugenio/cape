package primitives

import (
	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Service represents a service connecting to the system (e.g. a machine
// running a pipeline).
type Service struct {
	*database.Primitive
	Email       string            `json:"email"`
	Credentials *auth.Credentials `json:"credentials"`
}

// GetType returns the type for this entity
func (s *Service) GetType() types.Type {
	return ServiceType
}

// NewService returns a mutable service struct
func NewService() (*Service, error) {
	p, err := database.NewPrimitive(ServiceType)
	if err != nil {
		return nil, err
	}

	return &Service{
		Primitive: p,
	}, nil
}

// GetCredentials satisfies Identity interface
func (s *Service) GetCredentials() *auth.Credentials {
	return s.Credentials
}
