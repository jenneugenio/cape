package framework

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/primitives"
)

// IsAuthenticatedFn is a type alias for the function signaute of the directive
type IsAuthenticatedFn func(context.Context, interface{},
	graphql.Resolver, primitives.TokenType) (interface{}, error)

// IsAuthenticatedDirective checks to make sure a query is authenticated
func IsAuthenticatedDirective(db database.Backend, tokenAuthority *auth.TokenAuthority) IsAuthenticatedFn {
	return func(ctx context.Context, obj interface{},
		next graphql.Resolver, typeArg primitives.TokenType) (interface{}, error) {
		logger := Logger(ctx)
		token := AuthToken(ctx)

		err := tokenAuthority.Verify(token)
		if err != nil {
			msg := "Could not authenticates. Unable to verify auth token"
			logger.Info().Err(err).Msg(msg)
			return nil, ErrAuthentication
		}

		session := &primitives.Session{}
		err = db.QueryOne(ctx, session, database.NewFilter(database.Where{"token": token.String()}, nil, nil))
		if err != nil {
			msg := "Could not authenticates. Unable to find session"
			logger.Info().Err(err).Msg(msg)
			return nil, ErrAuthentication
		}

		user := &primitives.User{}
		err = db.Get(ctx, session.IdentityID, user)
		if err != nil {
			msg := "Could not authenticates. Unable to find user"
			logger.Error().Err(err).Msg(msg)
			return nil, ErrAuthentication
		}

		ctx = context.WithValue(ctx, SessionContextKey, session)
		ctx = context.WithValue(ctx, UserContextKey, user)

		return next(ctx)
	}
}
