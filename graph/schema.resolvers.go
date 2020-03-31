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

func (r *mutationResolver) Setup(ctx context.Context, input model.NewUserRequest) (*primitives.User, error) {
	// Make the user
	creds := &primitives.Credentials{
		PublicKey: &input.PublicKey,
		Salt:      &input.Salt,
		Alg:       input.Alg,
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

	// Make the role
	role, err := primitives.NewRole("admin", true)
	if err != nil {
		return nil, err
	}

	err = tx.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	// Assign the role
	assignment, err := primitives.NewAssignment(user.ID, role.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Create(ctx, assignment)
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
	serviceID := database.EmptyID
	if input.ServiceID != nil {
		service := &primitives.Service{}
		err := r.Backend.Get(ctx, *input.ServiceID, service)
		if err != nil {
			return nil, err
		}

		if service.Type != primitives.DataConnectorServiceType {
			return nil, errs.New(MustBeDataConnector, "Linking service to data source must be a data connector")
		}

		serviceID = *input.ServiceID
	}

	source, err := primitives.NewSource(input.Label, input.Credentials, serviceID)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *mutationResolver) RemoveSource(ctx context.Context, input model.RemoveSourceRequest) (*string, error) {
	var source primitives.Source
	filter := database.Filter{Where: database.Where{"label": input.Label}}
	err := r.Backend.QueryOne(ctx, &source, filter)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Delete(ctx, source.ID)
	return nil, err
}

func (r *mutationResolver) CreateLoginSession(ctx context.Context, input model.LoginSessionRequest) (*primitives.Session, error) {
	logger := fw.Logger(ctx)

	isFakeIdentity := false
	identity, err := queryIdentity(ctx, r.Backend, input.Email)
	if err != nil && !errs.FromCause(err, database.NotFoundCause) {
		logger.Info().Err(err).Msg("Could not authenticate. Error querying database")
		return nil, fw.ErrAuthentication
	} else if errs.FromCause(err, database.NotFoundCause) {
		// if identity doesn't exist need to return fake data
		isFakeIdentity = true
		identity, err = auth.NewFakeIdentity(input.Email)
		if err != nil {
			logger.Info().Err(err).Msg("Could not authenticate. Unable to create fake identity")
			return nil, fw.ErrAuthentication
		}
	}

	token, expiresIn, err := r.TokenAuthority.Generate(primitives.Login)
	if err != nil {
		logger.Info().Err(err).Msgf("Could not authenticate identity %s. Failed to generate auth token", identity.GetEmail())
		return nil, fw.ErrAuthentication
	}

	session, err := primitives.NewSession(identity, expiresIn, primitives.Login, token)
	if err != nil {
		logger.Info().Err(err).Msgf("Could not authenticate identity %s. Failed to create session", identity.GetEmail())
		return nil, fw.ErrAuthentication
	}

	if isFakeIdentity {
		// fake data doesn't need to be put in database so
		// return early
		return session, nil
	}

	err = r.Backend.Create(ctx, session)
	if err != nil {
		logger.Error().Err(err).Msgf("Could not authenticate identity %s. Create session in database", identity.GetEmail())
		return nil, fw.ErrAuthentication
	}

	return session, nil
}

func (r *mutationResolver) CreateAuthSession(ctx context.Context, input model.AuthSessionRequest) (*primitives.Session, error) {
	logger := fw.Logger(ctx)

	session := fw.Session(ctx)
	identity := fw.Identity(ctx)

	creds, err := auth.LoadCredentials(identity.GetCredentials().PublicKey, identity.GetCredentials().Salt)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Load credentials failed", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	err = creds.Verify(session.Token, &input.Signature)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Token verification failed", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	token, expiresIn, err := r.TokenAuthority.Generate(primitives.Authenticated)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Failed to generate auth token", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	authSession, err := primitives.NewSession(identity, expiresIn, primitives.Authenticated, token)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Failed to create session", identity.GetEmail())
		logger.Info().Err(err).Msg(msg)
		return nil, fw.ErrAuthentication
	}

	err = r.Backend.Create(ctx, authSession)
	if err != nil {
		msg := fmt.Sprintf("Could not authenticate identity %s. Create session in database", identity.GetEmail())
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

func (r *queryResolver) User(ctx context.Context, id database.ID) (*primitives.User, error) {
	return nil, errs.New(RouteNotImplemented, "User query not implemented")
}

func (r *queryResolver) Users(ctx context.Context) ([]*primitives.User, error) {
	return nil, errs.New(RouteNotImplemented, "Users query not implemented")
}

func (r *queryResolver) Me(ctx context.Context) (primitives.Identity, error) {
	identity := fw.Identity(ctx)
	return identity, nil
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

func (r *queryResolver) Identities(ctx context.Context, emails []*primitives.Email) ([]primitives.Identity, error) {
	serviceEmails := database.In{}
	userEmails := database.In{}

	for _, email := range emails {
		if email.Type == primitives.ServiceEmail {
			serviceEmails = append(serviceEmails, email.String())
		} else {
			userEmails = append(userEmails, email.String())
		}
	}

	users := []primitives.User{}
	if len(userEmails) > 0 {
		err := r.Backend.Query(ctx, &users, database.NewFilter(database.Where{"email": userEmails}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	services := []primitives.Service{}
	if len(serviceEmails) > 0 {
		err := r.Backend.Query(ctx, &services, database.NewFilter(database.Where{"email": serviceEmails}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	identities := make([]primitives.Identity, len(users)+len(services))
	for i, user := range users {
		u := &primitives.User{}
		*u = user
		identities[i] = u
	}

	for i, service := range services {
		s := &primitives.Service{}
		*s = service
		identities[i+len(users)] = s
	}

	return identities, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
