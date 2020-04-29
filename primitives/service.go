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
	Type     ServiceType `json:"type"`
	Endpoint *URL        `json:"endpoint"`
}

// GetType returns the type for this entity
func (s *Service) GetType() types.Type {
	return ServicePrimitiveType
}

// NewService returns a mutable service struct
func NewService(email Email, typ ServiceType, endpoint *URL, creds *Credentials) (*Service, error) {
	p, err := database.NewPrimitive(ServicePrimitiveType)
	if err != nil {
		return nil, err
	}

	service := &Service{
		IdentityImpl: &IdentityImpl{
			Primitive:   p,
			Email:       email,
			Credentials: creds,
		},
		Type:     typ,
		Endpoint: endpoint,
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
	case DataConnectorServiceType:
		if s.Endpoint == nil {
			return errors.New(InvalidServiceCause, "Must specify endpoint with data-connector service type")
		}

		if err := s.Endpoint.Validate(); err != nil {
			return err
		}

	case UserServiceType:
		if s.Endpoint != nil {
			return errors.New(InvalidServiceCause, "Can't specify endpoint on user service type")
		}

	default:
		return errors.New(InvalidServiceCause, "Unrecognized type: %s", s.Type.String())
	}

	if err := s.Credentials.Validate(); err != nil {
		return errors.Wrap(InvalidServiceCause, err)
	}

	return nil
}

// GetCredentials satisfies Identity interface
func (s *Service) GetCredentials() *Credentials {
	return s.Credentials
}

// GetEmail satisfies Identity interface
func (s *Service) GetEmail() Email {
	return s.Email
}
