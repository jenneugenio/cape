package primitives

import (
	"github.com/dropoutlabs/privacyai/database"
	"github.com/dropoutlabs/privacyai/database/types"
)

// Service represents a service connecting to the system (e.g. a machine
// running a pipeline).
type Service struct {
	*database.Primitive
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
