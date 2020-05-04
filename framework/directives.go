package framework

import (
	"context"

	"github.com/99designs/gqlgen/graphql"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
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
			return nil, auth.ErrAuthentication
		}

		err := tokenAuthority.Verify(token)
		if err != nil {
			msg := "Could not authenticate. Unable to verify auth token"
			logger.Info().Err(err).Msg(msg)
			return nil, auth.ErrAuthentication
		}

		session := &primitives.Session{}
		err = db.QueryOne(ctx, session, database.NewFilter(database.Where{"token": token.String()}, nil, nil))
		if err != nil {
			msg := "Could not authenticate. Unable to find session"
			logger.Info().Err(err).Msg(msg)
			return nil, auth.ErrAuthentication
		}

		ownerType, err := session.OwnerID.Type()
		if err != nil {
			msg := "Could not authenticate. Unable get credentialProvider type"
			logger.Info().Err(err).Msg(msg)
			return nil, auth.ErrAuthentication
		}

		identityType, err := session.IdentityID.Type()
		if err != nil {
			msg := "Could not authenticate. Unable get credentialProvider type"
			logger.Info().Err(err).Msg(msg)
			return nil, auth.ErrAuthentication
		}

		var credentialProvider primitives.CredentialProvider
		if ownerType == primitives.UserType {
			user := &primitives.User{}
			err = db.Get(ctx, session.IdentityID, user)
			if err != nil {
				msg := "Could not authenticate. Unable to find credentialProvider"
				logger.Error().Err(err).Msg(msg)
				return nil, auth.ErrAuthentication
			}
			credentialProvider = user
		} else if ownerType == primitives.TokenPrimitiveType {
			token := &primitives.Token{}
			err = db.Get(ctx, session.OwnerID, token)
			if err != nil {
				msg := "Could not authenticate. Unable to find credentialProvider"
				logger.Error().Err(err).Msg(msg)
				return nil, auth.ErrAuthentication
			}
			credentialProvider = token
		}

		var identity primitives.Identity
		if identityType == primitives.UserType {
			user := &primitives.User{}
			err = db.Get(ctx, session.IdentityID, user)
			if err != nil {
				msg := "Could not authenticate. Unable to find identity"
				logger.Error().Err(err).Msg(msg)
				return nil, auth.ErrAuthentication
			}
			identity = user
		} else if identityType == primitives.ServicePrimitiveType {
			service := &primitives.Service{}
			err = db.Get(ctx, session.IdentityID, service)
			if err != nil {
				msg := "Could not authenticate. Unable to find identity"
				logger.Error().Err(err).Msg(msg)
				return nil, auth.ErrAuthentication
			}
			identity = service
		}

		policies, err := QueryIdentityPolicies(ctx, db, identity.GetID())
		if err != nil {
			return nil, err
		}

		roles, err := QueryRoles(ctx, db, identity.GetID())
		if err != nil {
			return nil, err
		}

		aSession, err := auth.NewSession(identity, session, policies, roles, credentialProvider)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, SessionContextKey, aSession)

		return next(ctx)
	}
}
