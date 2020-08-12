package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/mailer"
)

// Resolver is used by graphql to resolve queries/mutations
type Resolver struct {
	Database           db.Interface
	CredentialProducer auth.CredentialProducer
	Mailer             mailer.Mailer
}
