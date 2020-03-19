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

func (r *mutationResolver) DeleteSession(ctx context.Context, input *model.DeleteSessionRequest) (string, error) {
	return "", errs.New(RouteNotImplemented, "Delete session not implemented")
}

func (r *queryResolver) User(ctx context.Context) (*primitives.User, error) {
	return nil, errs.New(RouteNotImplemented, "User query not implemented")
}

func (r *queryResolver) Session(ctx context.Context) (*primitives.Session, error) {
	return nil, errs.New(RouteNotImplemented, "Session query not implemented")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
