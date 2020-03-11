package primitives

import (
	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Identity represents an identity type such as user or service
type Identity interface {
	GetID() database.ID
	GetType() types.Type
	GetCredentials() *auth.Credentials
}
