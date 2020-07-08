package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserRequest) (*model.CreateUserResponse, error) {
	logger := fw.Logger(ctx)
	session := fw.Session(ctx)

	password, err := primitives.GeneratePassword()
	if err != nil {
		logger.Error().Err(err).Msg("Could not create password")
		return nil, err
	}

	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Info().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx) // nolint: errcheck

	enforcer := auth.NewEnforcer(session, tx)

	user, err := primitives.NewUser(input.Name, input.Email, creds)
	if err != nil {
		return nil, err
	}

	// We need to get the system roles back from the database so we can
	// assignment them to this user appropriately.
	systemRoles, err := fw.GetRolesByLabel(ctx, tx, []primitives.Label{
		primitives.GlobalRole,
	})
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	err = fw.CreateAssignments(ctx, tx, user, systemRoles)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &model.CreateUserResponse{
		Password: password,
		User:     user,
	}, nil
}

func (r *mutationResolver) DeleteSession(ctx context.Context, input model.DeleteSessionRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	if input.Token == nil {
		err := enforcer.Delete(ctx, primitives.SessionType, currSession.Session.ID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	found := false
	for _, role := range currSession.Roles {
		if role.Label == primitives.AdminRole {
			found = true
		}
	}

	if !found {
		return nil, errs.New(auth.AuthorizationFailure, "Unable to delete session")
	}

	id, err := r.TokenAuthority.Verify(input.Token)
	if err != nil {
		return nil, err
	}

	err = enforcer.Delete(ctx, primitives.SessionType, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *queryResolver) User(ctx context.Context, id database.ID) (*primitives.User, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	user := &primitives.User{}
	err := enforcer.Get(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*primitives.User, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var users []*primitives.User
	err := enforcer.Query(ctx, &users, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *queryResolver) Me(ctx context.Context) (*primitives.User, error) {
	session := fw.Session(ctx)
	return session.User, nil
}

func (r *queryResolver) Identities(ctx context.Context, emails []*primitives.Email) ([]*primitives.User, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var users []*primitives.User
	err := enforcer.Query(ctx, &users, database.NewFilter(database.Where{"email": emails}, nil, nil))
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userResolver) Roles(ctx context.Context, obj *primitives.User) ([]*primitives.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	return fw.QueryRoles(ctx, enforcer, obj.ID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
