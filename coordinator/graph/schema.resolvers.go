package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) Setup(ctx context.Context, input model.NewUserRequest) (*primitives.User, error) {
	// Make the user
	creds, err := primitives.NewCredentials(&input.PublicKey, &input.Salt)
	if err != nil {
		return nil, err
	}

	user, err := primitives.NewUser(input.Name, input.Email, creds)
	if err != nil {
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // nolint: errcheck

	err = tx.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	err = createSystemRoles(ctx, tx)
	if err != nil {
		return nil, err
	}

	err = attachDefaultPolicy(ctx, tx)
	if err != nil {
		return nil, err
	}

	roles, err := getRolesByLabel(ctx, tx, []primitives.Label{
		primitives.GlobalRole,
		primitives.AdminRole,
	})
	if err != nil {
		return nil, err
	}

	err = createAssignments(ctx, tx, user, roles)
	if err != nil {
		return nil, err
	}

	config, err := primitives.NewConfig()
	if err != nil {
		return nil, err
	}

	err = tx.Create(ctx, config)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUserRequest) (*primitives.User, error) {
	session := fw.Session(ctx)

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx) // nolint: errcheck

	enforcer := auth.NewEnforcer(session, tx)

	creds, err := primitives.NewCredentials(&input.PublicKey, &input.Salt)
	if err != nil {
		return nil, err
	}

	user, err := primitives.NewUser(input.Name, input.Email, creds)
	if err != nil {
		return nil, err
	}

	// We need to get the system roles back from the database so we can
	// assignment them to this user appropriately.
	systemRoles, err := getRolesByLabel(ctx, tx, []primitives.Label{
		primitives.GlobalRole,
	})
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	err = createAssignments(ctx, tx, user, systemRoles)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) AddSource(ctx context.Context, input model.AddSourceRequest) (*primitives.Source, error) {
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	if input.ServiceID != nil {
		service := &primitives.Service{}
		err := enforcer.Get(ctx, *input.ServiceID, service)
		if err != nil {
			return nil, err
		}

		if service.Type != primitives.DataConnectorServiceType {
			return nil, errs.New(MustBeDataConnector, "Linking service to data source must be a data connector")
		}
	}

	source, err := primitives.NewSource(input.Label, &input.Credentials, input.ServiceID)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *mutationResolver) RemoveSource(ctx context.Context, input model.RemoveSourceRequest) (*string, error) {
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	source := primitives.Source{}
	filter := database.Filter{Where: database.Where{"label": input.Label}}
	err := enforcer.QueryOne(ctx, &source, filter)
	if err != nil {
		return nil, err
	}

	err = enforcer.Delete(ctx, primitives.SourcePrimitiveType, source.ID)
	return nil, err
}

func (r *mutationResolver) CreateLoginSession(ctx context.Context, input model.LoginSessionRequest) (*primitives.Session, error) {
	logger := fw.Logger(ctx)

	isFakeIdentity := false
	identity, err := queryIdentity(ctx, r.Backend, input.Email)
	if err != nil && !errs.FromCause(err, database.NotFoundCause) {
		logger.Info().Err(err).Msg("Could not authenticate. Error querying database")
		return nil, auth.ErrAuthentication
	} else if errs.FromCause(err, database.NotFoundCause) {
		// if identity doesn't exist need to return fake data
		isFakeIdentity = true
		identity, err = auth.NewFakeIdentity(input.Email)
		if err != nil {
			logger.Info().Err(err).Msg("Could not authenticate. Unable to create fake identity")
			return nil, auth.ErrAuthentication
		}
	}

	token, expiresIn, err := r.TokenAuthority.Generate(primitives.Login)
	if err != nil {
		logger.Info().Err(err).Msgf("Could not authenticate identity %s. Failed to generate auth token", identity.GetEmail())
		return nil, auth.ErrAuthentication
	}

	session, err := primitives.NewSession(identity, expiresIn, primitives.Login, token)
	if err != nil {
		logger.Info().Err(err).Msgf("Could not authenticate identity %s. Failed to create session", identity.GetEmail())
		return nil, auth.ErrAuthentication
	}

	if isFakeIdentity {
		// fake data doesn't need to be put in database so
		// return early
		return session, nil
	}

	err = r.Backend.Create(ctx, session)
	if err != nil {
		logger.Error().Err(err).Msgf("Could not authenticate identity %s. Create session in database", identity.GetEmail())
		return nil, auth.ErrAuthentication
	}

	return session, nil
}

func (r *mutationResolver) CreateAuthSession(ctx context.Context, input model.AuthSessionRequest) (*primitives.Session, error) {
	logger := fw.Logger(ctx)
	s := fw.Session(ctx)

	enforcer := auth.NewEnforcer(s, r.Backend)

	session := s.Session
	identity := s.Identity

	creds, err := auth.LoadCredentials(identity.GetCredentials().PublicKey, identity.GetCredentials().Salt)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Load credentials failed", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, auth.ErrAuthentication
	}

	err = creds.Verify(session.Token, &input.Signature)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Token verification failed", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, auth.ErrAuthentication
	}

	token, expiresIn, err := r.TokenAuthority.Generate(primitives.Authenticated)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Failed to generate auth token", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, auth.ErrAuthentication
	}

	authSession, err := primitives.NewSession(identity, expiresIn, primitives.Authenticated, token)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Failed to create session", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, auth.ErrAuthentication
	}

	err = enforcer.Create(ctx, authSession)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Create session in database", identity.GetEmail())
		logger.Error().Err(err).Msg(msg)
		return nil, auth.ErrAuthentication
	}

	return authSession, nil
}

func (r *mutationResolver) DeleteSession(ctx context.Context, input model.DeleteSessionRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	session := &primitives.Session{}
	err := enforcer.QueryOne(ctx, session, database.NewFilter(database.Where{"token": input.Token.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	err = enforcer.Delete(ctx, primitives.SessionType, session.ID)
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
	return nil, errs.New(RouteNotImplemented, "Users query not implemented")
}

func (r *queryResolver) Me(ctx context.Context) (primitives.Identity, error) {
	session := fw.Session(ctx)
	return session.Identity, nil
}

func (r *queryResolver) Sources(ctx context.Context) ([]*primitives.Source, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var sources []*primitives.Source
	err := enforcer.Query(ctx, &sources, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func (r *queryResolver) Source(ctx context.Context, id database.ID) (*primitives.Source, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	source := &primitives.Source{}
	err := enforcer.Get(ctx, id, source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *queryResolver) SourceByLabel(ctx context.Context, label primitives.Label) (*primitives.Source, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	source := &primitives.Source{}
	err := enforcer.QueryOne(ctx, source, database.NewFilter(database.Where{"label": label.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *queryResolver) Identities(ctx context.Context, emails []*primitives.Email) ([]primitives.Identity, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	serviceEmails := database.In{}
	userEmails := database.In{}

	for _, email := range emails {
		if email.Type == primitives.ServiceEmail {
			serviceEmails = append(serviceEmails, email.String())
		} else {
			userEmails = append(userEmails, email.String())
		}
	}

	var users []*primitives.User
	if len(userEmails) > 0 {
		err := enforcer.Query(ctx, &users, database.NewFilter(database.Where{"email": userEmails}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	var services []*primitives.Service
	if len(serviceEmails) > 0 {
		err := enforcer.Query(ctx, &services, database.NewFilter(database.Where{"email": serviceEmails}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	identities := make([]primitives.Identity, len(users)+len(services))
	for i, user := range users {
		identities[i] = user
	}

	for i, service := range services {
		identities[i+len(users)] = service
	}

	return identities, nil
}

func (r *userResolver) Roles(ctx context.Context, obj *primitives.User) ([]*primitives.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	return getRoles(ctx, enforcer, obj.ID)
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
