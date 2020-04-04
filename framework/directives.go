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

		if token == nil {
			msg := "Could not authenticate. Token missing"
			logger.Info().Msg(msg)
			return nil, ErrAuthentication
		}

		err := tokenAuthority.Verify(token)
		if err != nil {
			msg := "Could not authenticate. Unable to verify auth token"
			logger.Info().Err(err).Msg(msg)
			return nil, ErrAuthentication
		}

		session := &primitives.Session{}
		err = db.QueryOne(ctx, session, database.NewFilter(database.Where{"token": token.String()}, nil, nil))
		if err != nil {
			msg := "Could not authenticate. Unable to find session"
			logger.Info().Err(err).Msg(msg)
			return nil, ErrAuthentication
		}

		typ, err := session.IdentityID.Type()
		if err != nil {
			msg := "Could not authenticate. Unable get identity type"
			logger.Info().Err(err).Msg(msg)
			return nil, ErrAuthentication
		}

		var identity primitives.Identity
		if typ == primitives.UserType {
			user := &primitives.User{}
			err = db.Get(ctx, session.IdentityID, user)
			if err != nil {
				msg := "Could not authenticate. Unable to find identity"
				logger.Error().Err(err).Msg(msg)
				return nil, ErrAuthentication
			}
			identity = user
		} else if typ == primitives.ServicePrimitiveType {
			service := &primitives.Service{}
			err = db.Get(ctx, session.IdentityID, service)
			if err != nil {
				msg := "Could not authenticate. Unable to find identity"
				logger.Error().Err(err).Msg(msg)
				return nil, ErrAuthentication
			}
			identity = service
		}

		ctx = context.WithValue(ctx, SessionContextKey, session)
		ctx = context.WithValue(ctx, IdentityContextKey, identity)

		return next(ctx)
	}
}
