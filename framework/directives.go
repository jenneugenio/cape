package framework

import (
	"context"

	"github.com/99designs/gqlgen/graphql"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/primitives"
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

		typ, err := session.OwnerID.Type()
		if err != nil {
			msg := "Could not authenticate. Unable get credentialProvider type"
			logger.Info().Err(err).Msg(msg)
			return nil, ErrAuthentication
		}

		var credentialProvider primitives.CredentialProvider
		if typ == primitives.UserType {
			user := &primitives.User{}
			err = db.Get(ctx, session.IdentityID, user)
			if err != nil {
				msg := "Could not authenticate. Unable to find credentialProvider"
				logger.Error().Err(err).Msg(msg)
				return nil, ErrAuthentication
			}
			credentialProvider = user
		} else if typ == primitives.TokenPrimitiveType {
			token := &primitives.TokenCredentials{}
			err = db.Get(ctx, session.OwnerID, token)
			if err != nil {
				msg := "Could not authenticate. Unable to find credentialProvider"
				logger.Error().Err(err).Msg(msg)
				return nil, ErrAuthentication
			}
			credentialProvider = token
		}

		ctx = context.WithValue(ctx, SessionContextKey, session)
		ctx = context.WithValue(ctx, CredentialProviderContextKey, credentialProvider)

		return next(ctx)
	}
}
