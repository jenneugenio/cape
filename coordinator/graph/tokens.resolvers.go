package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
	errs "github.com/capeprivacy/cape/partyerrors"
)

func (r *mutationResolver) CreateToken(ctx context.Context, input model.CreateTokenRequest) (*model.CreateTokenResponse, error) {
	logger := fw.Logger(ctx)
	currSession := fw.Session(ctx)
	user := currSession.User

	if user.ID != input.UserID && !currSession.Roles.Global.Can(models.CreateAnyToken) {
		return nil, errs.New(auth.AuthorizationFailure, "invalid permissions to create a token")
	} else if !currSession.Roles.Global.Can(models.CreateOwnToken) {
		return nil, errs.New(auth.AuthorizationFailure, "invalid permissions to create a token")
	}

	password := models.GeneratePassword()
	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Info().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	token := models.NewToken(input.UserID, &models.Credentials{
		Secret: creds.Secret,
		Salt:   creds.Salt,
		Alg:    creds.Alg,
	})

	err = r.Database.Tokens().Create(ctx, token)
	if err != nil {
		return nil, err
	}

	return &model.CreateTokenResponse{
		Secret: password,
		Token:  &token,
	}, nil
}

func (r *mutationResolver) RemoveToken(ctx context.Context, id string) (string, error) {
	currSession := fw.Session(ctx)

	if !currSession.Roles.Global.Can(models.RemoveAnyToken) {
		return "", errs.New(auth.AuthorizationFailure, "invalid permissions to remove a token")
	}

	err := r.Database.Tokens().Delete(ctx, id)
	if err != nil {
		return "", err
	}

	return id, err
}

func (r *queryResolver) Tokens(ctx context.Context, userID string) ([]string, error) {
	currSession := fw.Session(ctx)

	if currSession.User.ID != userID && !currSession.Roles.Global.Can(models.ListAnyTokens) {
		return nil, errs.New(auth.AuthorizationFailure, "unable to list tokens")
	} else if !currSession.Roles.Global.Can(models.ListOwnTokens) {
		return nil, errs.New(auth.AuthorizationFailure, "unable to list tokens")
	}

	tokens, err := r.Database.Tokens().ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tokenIDs := make([]string, len(tokens))
	for i, t := range tokens {
		tokenIDs[i] = t.ID
	}

	return tokenIDs, nil
}
