package graph

//go:generate go run github.com/99designs/gqlgen

import "github.com/dropoutlabs/cape/database"

type Resolver struct {
	Backend database.Backend
}
