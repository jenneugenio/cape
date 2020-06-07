package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateRecovery(ctx context.Context, input model.CreateRecoveryRequest) (*string, error) {
	logger := fw.Logger(ctx)

	if err := input.Email.Validate(); err != nil {
		logger.Info().Msg("Invalid email provided to create recovery")
		return nil, err
	}

	user := &primitives.User{}
	err := r.Backend.QueryOne(ctx, user, database.NewFilter(database.Where{
		"email": input.Email,
	}, nil, nil))
	if err != nil {
		// If the error is not found, we don't propagate it up, we pretend
		// everything is groovy so an attacker can't enumerate email addresses
		// through our recovery API
		if errs.FromCause(err, database.NotFoundCause) {
			logger.Info().Err(err).Msg("Could not find account to recover")
			return nil, nil
		}

		logger.Error().Err(err).Msg("Could not retrieve account for recovery")
		return nil, err
	}

	logger = logger.With().Str("user_id", user.ID.String()).Logger()

	password, err := primitives.GeneratePassword()
	if err != nil {
		logger.Error().Err(err).Msg("Could not generate password")
		return nil, err
	}

	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Error().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	recovery, err := primitives.NewRecovery(user.ID, creds)
	if err != nil {
		logger.Info().Err(err).Msg("Could not instantiate recovery")
		return nil, err
	}

	err = r.Backend.Create(ctx, recovery)
	if err != nil {
		logger.Error().Err(err).Msg("Could not insert recovery into database")
		return nil, err
	}

	err = r.mailer.Send(ctx, user, recovery, password)
	if err != nil {
		logger.Error().Err(err).Msg("Could not send recovery email to user")
		return nil, err
	}

	logger.Info().Msgf("Recovery created with id %s with secret %s", recovery.ID, password)
	return nil, nil
}

func (r *mutationResolver) AttemptRecovery(ctx context.Context, input model.AttemptRecoveryRequest) (*string, error) {
	logger := fw.Logger(ctx)

	if err := input.NewPassword.Validate(); err != nil {
		logger.Info().Err(err).Msg("Invalid password provided to attempt recovery")
		return nil, ErrRecoveryFailed
	}

	if err := input.Secret.Validate(); err != nil {
		logger.Info().Err(err).Msg("Invalid secret provided to attempt recovery")
		return nil, ErrRecoveryFailed
	}

	if err := input.ID.Validate(); err != nil {
		logger.Info().Err(err).Msg("Invalid id provided to attempt recovery")
		return nil, ErrRecoveryFailed
	}

	logger = logger.With().Str("recovery_id", input.ID.String()).Logger()

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not create transaction")
		return nil, ErrRecoveryFailed
	}
	defer tx.Rollback(ctx) // nolint: errcheck

	recovery := &primitives.Recovery{}
	err = tx.Get(ctx, input.ID, recovery)
	if err != nil {
		logger.Error().Err(err).Msg("Could not retrieve recovery")
		return nil, ErrRecoveryFailed
	}

	logger = logger.With().Str("user_id", recovery.UserID.String()).Logger()

	err = r.CredentialProducer.Compare(input.Secret, recovery.Credentials)
	if err != nil {
		logger.Info().Err(err).Msg("Invalid credentials provided")
		return nil, ErrRecoveryFailed
	}

	user := &primitives.User{}
	err = tx.Get(ctx, recovery.UserID, user)
	if err != nil {
		logger.Error().Err(err).Msg("Could not retrieve user for recovery")
		return nil, ErrRecoveryFailed
	}

	creds, err := r.CredentialProducer.Generate(input.NewPassword)
	if err != nil {
		logger.Error().Err(err).Msg("Could not generate creds for user")
		return nil, ErrRecoveryFailed
	}

	user.Credentials = creds
	err = tx.Update(ctx, user)
	if err != nil {
		logger.Error().Err(err).Msg("Could not update user with new password")
		return nil, ErrRecoveryFailed
	}

	err = tx.Delete(ctx, primitives.RecoveryType, recovery.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Could not ddelete recovery")
		return nil, ErrRecoveryFailed
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not commit transaction")
		return nil, ErrRecoveryFailed
	}

	logger.Info().Msg("Successfully recovered account with a new password")
	return nil, nil
}

func (r *mutationResolver) DeleteRecoveries(ctx context.Context, input model.DeleteRecoveriesRequest) (*string, error) {
	logger := fw.Logger(ctx)
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	// Only the worker can call this endpoint
	err := enforcer.Delete(ctx, primitives.RecoveryType, input.Ids...)
	if err != nil {
		logger.Error().Err(err).Msg("Could not delete recoveries")
		return nil, err
	}

	logger.Info().Msgf("Deleted %d recoveries", len(input.Ids))
	return nil, nil
}

func (r *queryResolver) Recoveries(ctx context.Context) ([]*primitives.Recovery, error) {
	logger := fw.Logger(ctx)
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	// TODO: Add ability for worker to filter recoveries so we only return the
	// recoveries we need to delete.
	recoveries := []*primitives.Recovery{}
	err := enforcer.Query(ctx, &recoveries, database.NewEmptyFilter())
	if err != nil {
		logger.Error().Err(err).Msg("Could not retrieve recoveries")
		return nil, err
	}

	logger.Info().Msgf("Retrieved %d recoveries", len(recoveries))
	return recoveries, nil
}
