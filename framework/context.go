package framework

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/gofrs/uuid"
	"github.com/manifoldco/go-base64"
	"github.com/rs/zerolog"
)

// ContextKey is a type alias used for storing data in a context
type ContextKey string

const (
	// RequestIDContextKey is the name of the key stored on the contet
	RequestIDContextKey ContextKey = "request-id"

	// LoggerContextKey is the name of the logger key stored on the context
	LoggerContextKey ContextKey = "logger"

	// AuthTokenContextKey is the name of the auth token stored on the context
	AuthTokenContextKey ContextKey = "auth-token"

	// SessionContextKey is the name of the session key stored on the context
	SessionContextKey ContextKey = "session"
)

// RequestID returns the request id stored on a given context
func RequestID(ctx context.Context) uuid.UUID {
	id := ctx.Value(RequestIDContextKey)
	if id == nil {
		panic("request id not available on context")
	}

	return id.(uuid.UUID)
}

// Logger returns the logger stored on a given context
func Logger(ctx context.Context) zerolog.Logger {
	logger := ctx.Value(LoggerContextKey)
	if logger == nil {
		panic("logger not available on context")
	}

	return logger.(zerolog.Logger)
}

// AuthToken returns the auth token stored on the given context.
//
// Returns nil if the token is not available on the context.
func AuthToken(ctx context.Context) *base64.Value {
	token := ctx.Value(AuthTokenContextKey)
	if token == nil {
		return nil
	}

	return token.(*base64.Value)
}

// Session returns the session stored on the given context
func Session(ctx context.Context) *auth.Session {
	session := ctx.Value(SessionContextKey)
	if session == nil {
		panic("session not available on context")
	}

	return session.(*auth.Session)
}
