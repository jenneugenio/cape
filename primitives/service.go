package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
	errors "github.com/dropoutlabs/cape/partyerrors"
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

	err = service.Validate()
	if err != nil {
		return nil, err
	}

	return service, nil
}

// Validate returns an error if the serive is not valid
func (s *Service) Validate() error {
	if s.Endpoint != nil && s.Type == UserServiceType {
		return errors.New(InvalidServiceCause, "Can't specify endpoint on user service type")
	}

	if s.Endpoint == nil && s.Type == DataConnectorServiceType {
		return errors.New(InvalidServiceCause, "Must specify endpoint with data-connector service type")
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
