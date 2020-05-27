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

func (r *mutationResolver) CreateToken(ctx context.Context, input model.CreateTokenRequest) (*model.CreateTokenResponse, error) {
	logger := fw.Logger(ctx)
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	identity := currSession.Identity

	hasAdminRole := hasRole(currSession.Roles, primitives.AdminRole)

	if hasAdminRole {
		// If the calling identity has the admin role they can only
		// create tokens for services and themselves. This should be
		// replaced by a policy DSL.

		typ, err := input.IdentityID.Type()
		if err != nil {
			return nil, err
		}

		if typ != primitives.ServicePrimitiveType && identity.GetID() != input.IdentityID {
			return nil, errs.New(auth.AuthorizationFailure,
				"An admin can only create tokens for a service and themselves")
		}
	} else if identity.GetID() != input.IdentityID {
		// If not an admin then the requesting identity must equal the requested identity

		return nil, errs.New(auth.AuthorizationFailure, "Can only create a token for yourself")
	}

	password, err := primitives.GeneratePassword()
	if err != nil {
		return nil, err
	}

	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Info().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	token, err := primitives.NewToken(input.IdentityID, creds)
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

	hasAdminRole := hasRole(currSession.Roles, primitives.AdminRole)

	token := &primitives.Token{}
	err := enforcer.Get(ctx, id, token)
	if err != nil {
		return database.EmptyID, err
	}

	if !hasAdminRole && currSession.Identity.GetID() != token.IdentityID {
		return database.EmptyID, errs.New(auth.AuthorizationFailure, "Can only remove tokens you own")
	}

	err = enforcer.Delete(ctx, primitives.TokenPrimitiveType, id)
	return id, err
}

func (r *queryResolver) Tokens(ctx context.Context, identityID database.ID) ([]database.ID, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	hasAdminRole := hasRole(currSession.Roles, primitives.AdminRole)

	if !hasAdminRole && currSession.Identity.GetID() != identityID {
		return nil, errs.New(auth.AuthorizationFailure, "Can only list tokens you own")
	}

	var tokens []*primitives.Token
	filter := database.Filter{Where: database.Where{"identity_id": identityID}}

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
