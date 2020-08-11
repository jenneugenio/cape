package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
)

func (r *mutationResolver) CreateRecovery(ctx context.Context, input model.CreateRecoveryRequest) (*string, error) {
	logger := fw.Logger(ctx)

	user, err := r.Database.Users().Get(ctx, input.Email)
	if err != nil {
		// If the error is not found, we don't propagate it up, we pretend
		// everything is groovy so an attacker can't enumerate email addresses
		// through our recovery API
		if err == db.ErrCannotFindUser {
			logger.Info().Err(err).Msg("Could not find account to recover")
			return nil, nil
		}

		logger.Error().Err(err).Msg("Could not retrieve account for recovery")
		return nil, err
	}

	logger = logger.With().Str("user_id", user.ID).Logger()

	password := models.GeneratePassword()

	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Error().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	recovery := models.NewRecovery(user.ID, &models.Credentials{
		Secret: creds.Secret,
		Salt:   creds.Salt,
		Alg:    creds.Alg,
	})

	err = r.Database.Recoveries().Create(ctx, recovery)
	if err != nil {
		logger.Error().Err(err).Msg("Could not insert recovery into database")
		return nil, err
	}

	err = r.Mailer.SendAccountRecovery(ctx, *user, recovery, password)
	if err != nil {
		logger.Error().Err(err).Msg("Could not send recovery email to user")
		return nil, err
	}

	logger.Info().Msgf("Recovery created with id %s with secret %s", recovery.ID, password)
	return nil, nil
}

func (r *mutationResolver) AttemptRecovery(ctx context.Context, input model.AttemptRecoveryRequest) (*string, error) {
	logger := fw.Logger(ctx)

	logger = logger.With().Str("recovery_id", input.ID).Logger()

	recovery, err := r.Database.Recoveries().Get(ctx, input.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Could not retrieve recovery")
		return nil, ErrRecoveryFailed
	}

	logger = logger.With().Str("user_id", recovery.UserID).Logger()

	if recovery.Expired() {
		logger.Info().Msg("Recovery has expired")
		return nil, ErrRecoveryFailed
	}

	err = r.CredentialProducer.Compare(input.Secret, recovery.Credentials)
	if err != nil {
		logger.Info().Err(err).Msg("Invalid credentials provided")
		return nil, ErrRecoveryFailed
	}

	user, err := r.Database.Users().GetByID(ctx, recovery.UserID)
	if err != nil {
		logger.Error().Err(err).Msgf("Could not get user %s", recovery.UserID)
		return nil, ErrRecoveryFailed
	}

	creds, err := r.CredentialProducer.Generate(input.NewPassword)
	if err != nil {
		logger.Error().Err(err).Msg("Could not generate creds for user")
		return nil, ErrRecoveryFailed
	}

	user.Credentials = *creds
	err = r.Database.Users().Update(ctx, user.ID, *user)
	if err != nil {
		logger.Error().Err(err).Msg("Could not update user with new password")
		return nil, ErrRecoveryFailed
	}

	err = r.Database.Recoveries().Delete(ctx, recovery.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Could not ddelete recovery")
		return nil, ErrRecoveryFailed
	}

	logger.Info().Msg("Successfully recovered account with a new password")
	return nil, nil
}

func (r *recoveryResolver) CreatedAt(ctx context.Context, obj *models.Recovery) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *recoveryResolver) UpdatedAt(ctx context.Context, obj *models.Recovery) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

// Recovery returns generated.RecoveryResolver implementation.
func (r *Resolver) Recovery() generated.RecoveryResolver { return &recoveryResolver{r} }

type recoveryResolver struct{ *Resolver }
