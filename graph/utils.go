package graph

import (
	"context"

	"github.com/dropoutlabs/cape/database"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

func queryIdentity(ctx context.Context, db database.Backend, email string) (primitives.Identity, error) {
	filter := database.NewFilter(database.Where{"email": email}, nil, nil)

	user := &primitives.User{}
	err := db.QueryOne(ctx, user, filter)
	if err != nil && !errors.FromCause(err, database.NotFoundCause) {
		return nil, err
	}
	if err == nil {
		return user, err
	}

	service := &primitives.Service{}
	err = db.QueryOne(ctx, service, filter)
	if err != nil {
		return nil, err
	}

	return service, nil
}
