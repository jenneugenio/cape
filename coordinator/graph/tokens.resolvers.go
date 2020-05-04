package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateToken(ctx context.Context, input model.CreateTokenRequest) (*primitives.Token, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	creds, err := primitives.NewCredentials(&input.PublicKey, &input.Salt)
	if err != nil {
		return nil, err
	}

	token, err := primitives.NewToken(input.IdentityID, creds)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *mutationResolver) RemoveToken(ctx context.Context, id database.ID) (database.ID, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	err := enforcer.Delete(ctx, primitives.TokenPrimitiveType, id)
	return id, err
}

func (r *queryResolver) Tokens(ctx context.Context, identityID database.ID) ([]database.ID, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

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
