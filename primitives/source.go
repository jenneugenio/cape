package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
	"net/url"
)

// Source represents the connection information for an external data source
type Source struct {
	*database.Primitive
	Label    Label   `json:"label"`
	Endpoint url.URL `json:"endpoint"`

	Credentials url.URL `json:"credentials"`
}

// GetType returns the type for this entity
func (s *Source) GetType() types.Type {
	return SourceType
}

// NewSource returns a new Source struct
func NewSource(label Label, credentials url.URL) (*Source, error) {
	p, err := database.NewPrimitive(SourceType)
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
		Credentials: credentials,
		Endpoint:    *credentialCopy,
	}, nil
}
