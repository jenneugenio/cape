package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/primitives"
	"github.com/manifoldco/go-base64"
)

func (r *mutationResolver) CreateToken(ctx context.Context, input model.CreateTokenRequest) (*primitives.Token, error) {
	creds, err := primitives.NewCredentials(&input.PublicKey, &input.Salt)
	if err != nil {
		return nil, err
	}

	token, err := primitives.NewTokenCredentials(input.IdentityID, creds)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *tokenResolver) PublicKey(ctx context.Context, obj *primitives.Token) (base64.Value, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tokenResolver) Salt(ctx context.Context, obj *primitives.Token) (base64.Value, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tokenResolver) Alg(ctx context.Context, obj *primitives.Token) (primitives.CredentialsAlgType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Token returns generated.TokenResolver implementation.
func (r *Resolver) Token() generated.TokenResolver { return &tokenResolver{r} }

type tokenResolver struct{ *Resolver }
