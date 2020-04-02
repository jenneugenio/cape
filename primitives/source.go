package primitives

import (
	"net/url"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Source represents the connection information for an external data source
type Source struct {
	*database.Primitive
	Label       Label      `json:"label"`
	Endpoint    url.URL    `json:"endpoint"`
	Type        SourceType `json:"type"`
	Credentials url.URL    `json:"credentials"`

	// ServiceID can be nil as it's not set when a data connector has not been
	// linked with the service.
	ServiceID *database.ID `json:"service_id"`
}

// GetType returns the type for this entity
func (s *Source) GetType() types.Type {
	return SourcePrimitiveType
}

// Validate returns whether or not the source represents a valid Source
func (s *Source) Validate() error {
	if err := s.Primitive.Validate(); err != nil {
		return err
	}

	if err := s.Label.Validate(); err != nil {
		return err
	}

	if err := s.Type.Validate(); err != nil {
		return err
	}

	if s.ServiceID == nil {
		return nil
	}

	return s.ServiceID.Validate()
}

// NewSource returns a new Source struct
func NewSource(label Label, credentials url.URL, serviceID *database.ID) (*Source, error) {
	p, err := database.NewPrimitive(SourcePrimitiveType)
	if err != nil {
		return nil, err
	}

	t, err := NewSourceType(credentials.Scheme)
	if err != nil {
		return nil, err
	}

	credentialCopy, err := url.Parse(credentials.String())
	if err != nil {
		return nil, err
	}

	// delete the credential part of the URL
	credentialCopy.User = nil

	return &Source{
		Primitive:   p,
		Label:       label,
		Type:        t,
		Credentials: credentials,
		Endpoint:    *credentialCopy,
		ServiceID:   serviceID,
	}, nil
}
