package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/mailer"
	"github.com/capeprivacy/cape/coordinator/db"
)

// Resolver is used by graphql to resolve queries/mutations
type Resolver struct {
	Database           db.Interface
	Backend            database.Backend
	TokenAuthority     *auth.TokenAuthority
	CredentialProducer auth.CredentialProducer
	RootKey            [32]byte
	Mailer             mailer.Mailer
}
