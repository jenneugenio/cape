package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
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
func NewService(email Email, typ ServiceType) (*Service, error) {
	p, err := database.NewPrimitive(ServicePrimitiveType)
	if err != nil {
		return nil, err
	}

	name := Name(email.String())
	service := &Service{
		IdentityImpl: &IdentityImpl{
			Primitive: p,
			Email:     email,
			Name:      name,
		},
		Type: typ,
	}

	return service, service.Validate()
}

// Validate returns an error if the serive is not valid
func (s *Service) Validate() error {
	if err := s.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidServiceCause, err)
	}

	if err := s.Type.Validate(); err != nil {
		return errors.Wrap(InvalidServiceCause, err)
	}

	switch s.Type {
	case UserServiceType:
		return nil
	default:
		return errors.New(InvalidServiceCause, "Unrecognized type: %s", s.Type.String())
	}
}

// GetEmail satisfies Identity interface
func (s *Service) GetEmail() Email {
	return s.Email
}

func (s *Service) GetEncryptable() bool {
	return false
}
