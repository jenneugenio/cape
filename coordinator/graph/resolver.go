package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/database"
)

// Resolver is used by graphql to resolve queries/mutations
type Resolver struct {
	Backend        database.Backend
	TokenAuthority *auth.TokenAuthority
}
