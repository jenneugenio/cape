package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateToken(ctx context.Context, input model.CreateTokenRequest) (*primitives.Token, error) {
	creds, err := primitives.NewCredentials(&input.PublicKey, &input.Salt)
	if err != nil {
		return nil, err
	}

	token, err := primitives.NewTokenCredentials(input.IdentityID, creds)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
