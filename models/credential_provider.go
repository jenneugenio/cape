package models

type CredentialProvider interface {
	GetStringID() string
	GetCredentials() (*Credentials, error)
	GetUserID() string
}
