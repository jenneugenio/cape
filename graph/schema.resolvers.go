package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dropoutlabs/cape/graph/generated"
	"github.com/dropoutlabs/cape/graph/model"
	"github.com/dropoutlabs/cape/primitives"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUserRequest) (*primitives.User, error) {
	user, err := primitives.NewUser(input.Name, "", nil)
	if err != nil {
		return nil, err
	}

	user.ID = input.ID
	err = r.Backend.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queryResolver) User(ctx context.Context) (*primitives.User, error) {
	return &primitives.User{}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
