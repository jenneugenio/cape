package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/database"
	fw "github.com/dropoutlabs/cape/framework"
	"github.com/dropoutlabs/cape/graph/generated"
	"github.com/dropoutlabs/cape/graph/model"
	errs "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

func (r *assignmentResolver) Role(ctx context.Context, obj *primitives.Assignment) (*primitives.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *assignmentResolver) Identity(ctx context.Context, obj *primitives.Assignment) (primitives.Identity, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUserRequest) (*primitives.User, error) {
	creds := &primitives.Credentials{
		PublicKey: &input.PublicKey,
		Salt:      &input.Salt,
		Alg:       input.Alg,
	}

	user, err := primitives.NewUser(input.Name, input.Email, creds)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) AddSource(ctx context.Context, input model.AddSourceRequest) (*primitives.Source, error) {
	source, err := primitives.NewSource(input.Label, input.Credentials)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *mutationResolver) CreateLoginSession(ctx context.Context, input model.LoginSessionRequest) (*primitives.Session, error) {
	logger := fw.Logger(ctx)

	isFakeUser := false
	user := &primitives.User{}
	err := r.Backend.QueryOne(ctx, user, database.NewFilter(database.Where{"email": input.Email}, nil, nil))
	if err != nil {
		if errs.FromCause(err, database.NotFoundCause) {
			isFakeUser = true
			user, err = auth.NewFakeUser(input.Email)
			if err != nil {
				logger.Info().Err(err).Msg("Could not authenticate. Unable to create fake user")
				return nil, fw.ErrAuthentication
			}
		} else {
			logger.Info().Err(err).Msg("Could not authenticate. Error querying database")
			return nil, fw.ErrAuthentication
		}
	}

	token, expiresIn, err := r.TokenAuthority.Generate(primitives.Login)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Failed to generate auth token", user.Email)
		logger.Info().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	session, err := primitives.NewSession(user, expiresIn, primitives.Login, token)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Failed to create session", user.Email)
		logger.Info().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	if isFakeUser {
		return session, nil
	}

	err = r.Backend.Create(ctx, session)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Create session in database", user.Email)
		logger.Error().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	return session, nil
}

func (r *mutationResolver) CreateAuthSession(ctx context.Context, input model.AuthSessionRequest) (*primitives.Session, error) {
	logger := fw.Logger(ctx)

	session := fw.Session(ctx)
	user := fw.User(ctx)

	creds, err := auth.LoadCredentials(user.Credentials.PublicKey, user.Credentials.Salt)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Load credentials failed", user.Email)
		logger.Error().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	err = creds.Verify(session.Token, &input.Signature)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Token verification failed", user.Email)
		logger.Error().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	token, expiresIn, err := r.TokenAuthority.Generate(primitives.Authenticated)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Failed to generate auth token", user.Email)
		logger.Info().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	authSession, err := primitives.NewSession(user, expiresIn, primitives.Authenticated, token)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Failed to create session", user.Email)
		logger.Error().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	err = r.Backend.Create(ctx, authSession)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate user %s. Create session in database", user.Email)
		logger.Error().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	return authSession, nil
}

func (r *mutationResolver) DeleteSession(ctx context.Context, input model.DeleteSessionRequest) (*string, error) {
	session := &primitives.Session{}
	err := r.Backend.QueryOne(ctx, session, database.NewFilter(database.Where{"token": input.Token.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	err = r.Backend.Delete(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mutationResolver) CreateRole(ctx context.Context, input model.CreateRoleRequest) (*primitives.Role, error) {
	role, err := primitives.NewRole(input.Label)
	if err != nil {
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // nolint: errcheck

	err = tx.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	if len(input.IdentityIds) > 0 {
		assignments := make([]database.Entity, len(input.IdentityIds))
		for i, id := range input.IdentityIds {
			assignment, err := primitives.NewAssignment(id, role.ID)
			if err != nil {
				return nil, err
			}
			assignments[i] = assignment
		}
		err = tx.Create(ctx, assignments...)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *mutationResolver) DeleteRole(ctx context.Context, input model.DeleteRoleRequest) (*string, error) {
	err := r.Backend.Delete(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mutationResolver) AssignRole(ctx context.Context, input model.AssignRoleRequest) (*primitives.Assignment, error) {
	assignment, err := primitives.NewAssignment(input.AssignmentID, input.RoleID)
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (r *mutationResolver) UnassignRole(ctx context.Context, input model.AssignRoleRequest) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context) (*primitives.User, error) {
	return nil, errs.New(RouteNotImplemented, "User query not implemented")
}

func (r *queryResolver) Session(ctx context.Context) (*primitives.Session, error) {
	return nil, errs.New(RouteNotImplemented, "Session query not implemented")
}

func (r *queryResolver) Sources(ctx context.Context) ([]*primitives.Source, error) {
	var s []primitives.Source
	err := r.Backend.Query(ctx, &s, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	sources := make([]*primitives.Source, len(s))
	for i := 0; i < len(sources); i++ {
		sources[i] = &(s[i])
	}

	return sources, nil
}

func (r *queryResolver) Source(ctx context.Context, id database.ID) (*primitives.Source, error) {
	var primitive primitives.Source
	err := r.Backend.Get(ctx, id, &primitive)
	if err != nil {
		return nil, err
	}

	return &primitive, nil
}

func (r *queryResolver) Role(ctx context.Context, id database.ID) (*primitives.Role, error) {
	var primitive primitives.Role
	err := r.Backend.Get(ctx, id, &primitive)
	if err != nil {
		return nil, err
	}

	return &primitive, nil
}

func (r *queryResolver) Roles(ctx context.Context) ([]*primitives.Role, error) {
	var tmpR []primitives.Role
	err := r.Backend.Query(ctx, &tmpR, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	roles := make([]*primitives.Role, len(tmpR))
	for i := 0; i < len(roles); i++ {
		roles[i] = &(tmpR[i])
	}

	return roles, nil
}

func (r *queryResolver) RoleMembers(ctx context.Context, roleID database.ID) ([]primitives.Identity, error) {
	var a []primitives.Assignment
	err := r.Backend.Query(ctx, &a, database.NewFilter(database.Where{"role_id": roleID.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	userIDs := database.In{}
	serviceIDs := database.In{}
	for _, assignment := range a {
		typ, err := assignment.IdentityID.Type()
		if err != nil {
			return nil, err
		}

		if typ == primitives.UserType {
			userIDs = append(userIDs, assignment.IdentityID)
		} else if typ == primitives.ServiceType {
			serviceIDs = append(serviceIDs, assignment.IdentityID)
		}
	}

	var users []primitives.User
	if len(userIDs) > 0 {
		err = r.Backend.Query(ctx, &users, database.NewFilter(database.Where{"id": userIDs}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	var services []primitives.Service
	if len(serviceIDs) > 0 {
		err = r.Backend.Query(ctx, &services, database.NewFilter(database.Where{"id": serviceIDs}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	identities := make([]primitives.Identity, len(a))
	for i, user := range users {
		identities[i] = &user
	}

	for i, service := range services {
		identities[i+len(users)] = &service
	}

	return identities, nil
}

// Assignment returns generated.AssignmentResolver implementation.
func (r *Resolver) Assignment() generated.AssignmentResolver { return &assignmentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type assignmentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
