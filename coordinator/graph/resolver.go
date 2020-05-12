package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
)

// Resolver is used by graphql to resolve queries/mutations
type Resolver struct {
	Backend        database.Backend
	TokenAuthority *auth.TokenAuthority
	RootKey        [32]byte
}
