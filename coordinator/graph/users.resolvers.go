package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserRequest) (*model.CreateUserResponse, error) {
	logger := fw.Logger(ctx)
	session := fw.Session(ctx)

	password := primitives.GeneratePassword()

	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Info().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	user := models.NewUser(input.Name, input.Email, creds)

	err = r.Database.Users().Create(ctx, user)
	if err != nil {
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx) // nolint: errcheck

	enforcer := auth.NewEnforcer(session, tx)

	// We need to get the system roles back from the database so we can
	// assignment them to this user appropriately.
	systemRoles, err := fw.GetRolesByLabel(ctx, tx, []primitives.Label{
		primitives.GlobalRole,
	})
	if err != nil {
		return nil, err
	}

	err = fw.CreateAssignments(ctx, enforcer, user.ID, systemRoles)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
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

func (r *userResolver) Roles(ctx context.Context, obj *models.User) ([]*models.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	return fw.QueryRoles(ctx, enforcer, obj.ID)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
