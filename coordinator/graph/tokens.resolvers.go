package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/capeprivacy/cape/models"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateToken(ctx context.Context, input model.CreateTokenRequest) (*model.CreateTokenResponse, error) {
	logger := fw.Logger(ctx)
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	user := currSession.User

	if user.ID != input.UserID {
		return nil, errs.New(auth.AuthorizationFailure, "Can only create a token for yourself")
	}

	password := primitives.GeneratePassword()
	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Info().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	token, err := primitives.NewToken(input.UserID, &primitives.Credentials{
		Secret: creds.Secret,
		Salt:   creds.Salt,
		Alg:    primitives.CredentialsAlgType(creds.Alg),
	})
	if err != nil {
		logger.Info().Err(err).Msg("Could not create token")
		return nil, err
	}

	err = enforcer.Create(ctx, token)
	if err != nil {
		logger.Info().Err(err).Msg("Could not insert token into database")
		return nil, err
	}

	return &model.CreateTokenResponse{
		Secret: password,
		Token:  token,
	}, nil
}

func (r *mutationResolver) RemoveToken(ctx context.Context, id database.ID) (database.ID, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	hasAdminRole := hasRole(currSession.Roles, models.AdminRole)

	token := &primitives.Token{}
	err := enforcer.Get(ctx, id, token)
	if err != nil {
		return database.EmptyID, err
	}

	if !hasAdminRole && currSession.User.ID != token.UserID {
		return database.EmptyID, errs.New(auth.AuthorizationFailure, "Can only remove tokens you own")
	}

	err = enforcer.Delete(ctx, primitives.TokenPrimitiveType, id)
	return id, err
}

func (r *queryResolver) Tokens(ctx context.Context, userID string) ([]database.ID, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	hasAdminRole := hasRole(currSession.Roles, models.AdminRole)

	if !hasAdminRole && currSession.User.ID != userID {
		return nil, errs.New(auth.AuthorizationFailure, "Can only list tokens you own")
	}

	var tokens []*primitives.Token
	filter := database.Filter{Where: database.Where{"user_id": userID}}

	err := enforcer.Query(ctx, &tokens, filter)
	if err != nil {
		return nil, err
	}

	ids := make([]database.ID, len(tokens))
	for i, t := range tokens {
		ids[i] = t.ID
	}

	return ids, nil
}
