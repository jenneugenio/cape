package primitives

import "github.com/capeprivacy/cape/database"

type CredentialProvider interface {
	database.Entity
	GetCredentials() (*Credentials, error)
	GetIdentityID() database.ID
}
