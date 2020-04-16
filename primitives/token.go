package primitives

import (
	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/database/types"
)

// Token for an authorized entity (user or service)
type Token struct {
	*database.Primitive
	IdentityID database.ID `json:"identity_id"`
}

// GetType returns the type for this entity
func (t *Token) GetType() types.Type {
	return TokenPrimitiveType
}

// NewToken returns an immutable token struct
func NewToken(identityID database.ID) (*Token, error) {
	p, err := database.NewPrimitive(TokenPrimitiveType)
	if err != nil {
		return nil, err
	}

	// TODO: Figure out what fields should be set on this struct
	t := &Token{
		Primitive:  p,
		IdentityID: identityID,
	}

	ID, err := database.DeriveID(t)
	if err != nil {
		return nil, err
	}

	t.ID = ID
	return t, nil
}
