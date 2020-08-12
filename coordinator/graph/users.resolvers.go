package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserRequest) (*model.CreateUserResponse, error) {
	logger := fw.Logger(ctx)
	password := models.GeneratePassword()

	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Info().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	user := models.NewUser(input.Name, input.Email, *creds)

	err = r.Database.Users().Create(ctx, user)
	if err != nil {
		return nil, err
	}

	_, err = r.Database.Roles().SetOrgRole(ctx, user.Email, models.UserRole)
	if err != nil {
		return nil, err
	}

	return &model.CreateUserResponse{
		Password: password,
		User:     &user,
	}, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	return r.Database.Users().GetByID(ctx, id)
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users, err := r.Database.Users().List(ctx, nil)
	if err != nil {
		return nil, err
	}

	userPtrs := make([]*models.User, len(users))
	for i, user := range users {
		u := user
		userPtrs[i] = &u
	}
	return userPtrs, nil
}

func (r *userResolver) Role(ctx context.Context, obj *models.User) (*models.Role, error) {
	return r.Database.Roles().GetOrgRole(ctx, obj.Email)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
