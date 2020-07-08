package primitives

import "github.com/capeprivacy/cape/coordinator/database"

type CredentialProvider interface {
	database.Entity
	GetCredentials() (*Credentials, error)
	GetUserID() database.ID
}
