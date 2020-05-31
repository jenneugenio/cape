package primitives

import (
	"context"
	"encoding/json"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Source represents the connection information for an external data source
type Source struct {
	*database.Primitive
	Label Label      `json:"label"`
	Type  SourceType `json:"type"`

	// Endpoint is a "safe" version of the credential containing a hostname or
	// identifier for the underlying credential
	Endpoint *DBURL `json:"endpoint"`

	// XXX: Credentials contains a secret (user and password); it should only
	// _ever_ be returned to data connectors.
	Credentials *DBURL `json:"credentials"`

	// ServiceID can be nil as it's not set when a data connector has not been
	// linked with the service.
	ServiceID *database.ID `json:"service_id"`
}

type encryptedSource struct {
	*Source
	Credentials *base64.Value `json:"credentials"`
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

	if err := s.Endpoint.Validate(); err != nil {
		return err
	}

	if err := s.Credentials.Validate(); err != nil {
		return err
	}

	if s.ServiceID != nil {
		if err := s.ServiceID.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Encrypt implements the Encryptable interface
func (s *Source) Encrypt(ctx context.Context, codec crypto.EncryptionCodec) ([]byte, error) {
	data, err := codec.Encrypt(ctx, base64.New([]byte(s.Credentials.String())))
	if err != nil {
		return nil, err
	}

	return json.Marshal(encryptedSource{
		Source:      s,
		Credentials: data,
	})
}

// Decrypt implements the Encryptable interface
func (s *Source) Decrypt(ctx context.Context, codec crypto.EncryptionCodec, data []byte) error {
	in := &encryptedSource{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	unencrypted, err := codec.Decrypt(ctx, in.Credentials)
	if err != nil {
		return err
	}

	s.Primitive = in.Primitive
	s.Label = in.Label
	s.Type = in.Type
	s.Endpoint = in.Endpoint
	s.ServiceID = in.ServiceID

	u, err := NewDBURL(string(*unencrypted))
	if err != nil {
		return err
	}

	s.Credentials = u
	return nil
}

// NewSource returns a new Source struct
func NewSource(label Label, credentials *DBURL, serviceID *database.ID) (*Source, error) {
	p, err := database.NewPrimitive(SourcePrimitiveType)
	if err != nil {
		return nil, err
	}

	if credentials == nil {
		return nil, errors.New(InvalidSourceCause, "Credentials is a required field for a source")
	}

	t, err := NewSourceType(credentials.Scheme)
	if err != nil {
		return nil, err
	}

	endpoint, err := credentials.Copy()
	if err != nil {
		return nil, err
	}

	// delete the credential part of the URL for usage as the endpoint value.
	endpoint.User = nil

	source := &Source{
		Primitive:   p,
		Label:       label,
		Type:        t,
		Credentials: credentials,
		Endpoint:    endpoint,
		ServiceID:   serviceID,
	}

	return source, source.Validate()
}

func (s *Source) GetEncryptable() bool {
	return true
}
