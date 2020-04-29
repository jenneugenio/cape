package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/felixge/httpsnoop"
	"github.com/gofrs/uuid"
	"github.com/manifoldco/go-base64"
	"github.com/rs/zerolog"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
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

	// IdentityContextKey is the name of the identity key stored on the context
	IdentityContextKey ContextKey = "identity"
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
func Session(ctx context.Context) *primitives.Session {
	session := ctx.Value(SessionContextKey)
	if session == nil {
		panic("session not available on context")
	}

	return session.(*primitives.Session)
}

// Identity returns the identity stored on the given context
func Identity(ctx context.Context) primitives.Identity {
	identity := ctx.Value(IdentityContextKey)
	if identity == nil {
		panic("identity not available on context")
	}

	return identity.(primitives.Identity)
}

// RequestIDMiddleware sets a UUID on the response header and request context
// for use in tracing and
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		id, err := uuid.NewV4()
		if err != nil {
			panic(fmt.Sprintf("Could not generate a v4 uuid: %s", err))
		}

		ctx := context.WithValue(req.Context(), RequestIDContextKey, id)
		req = req.WithContext(ctx)
		rw.Header().Set("X-Request-ID", id.String())
		next.ServeHTTP(rw, req)
	})
}

// LogMiddleware sets a zerolog.Logger on the request context for use in
// downstream callers. This middleware relies on the requestIDMiddleware.
func LogMiddleware(log *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			requestID := RequestID(ctx)
			logger := log.With().Str("request_id", requestID.String()).Logger()

			ctx = context.WithValue(ctx, LoggerContextKey, logger)
			req = req.WithContext(ctx)
			next.ServeHTTP(rw, req)
		})
	}
}

// AuthTokenMiddleware sets the session ID on the request context for us in
// graphql handlers and elsewhere
func AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(rw, req)
			return
		}

		token, err := auth.GetBearerToken(authHeader)
		if err != nil {
			respondWithGQLError(rw, err)
			return
		}

		ctx = context.WithValue(ctx, AuthTokenContextKey, token)
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// IsAuthenticatedMiddleware sets the session ID on the request context for us in
// graphql handlers and elsewhere
func IsAuthenticatedMiddleware(db database.Backend, tokenAuthority *auth.TokenAuthority) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			route, err := getGraphqlRoute(req)
			if err != nil {
				fmt.Println("HELLO", err)
				respondWithGQLError(rw, ErrAuthentication)
				return
			}

			if route == "createLoginSession" {
				next.ServeHTTP(rw, req)
				return
			}

			logger := Logger(ctx)
			token := AuthToken(ctx)

			if token == nil {
				msg := "Could not authenticate. Token missing"
				logger.Info().Msg(msg)
				respondWithGQLError(rw, ErrAuthentication)
				return
			}

			err = tokenAuthority.Verify(token)
			if err != nil {
				msg := "Could not authenticate. Unable to verify auth token"
				logger.Info().Err(err).Msg(msg)
				respondWithGQLError(rw, ErrAuthentication)
				return
			}

			session := &primitives.Session{}
			err = db.QueryOne(ctx, session, database.NewFilter(database.Where{"token": token.String()}, nil, nil))
			if err != nil {
				msg := "Could not authenticate. Unable to find session"
				logger.Info().Err(err).Msg(msg)
				respondWithGQLError(rw, ErrAuthentication)
				return
			}

			typ, err := session.IdentityID.Type()
			if err != nil {
				msg := "Could not authenticate. Unable get identity type"
				logger.Info().Err(err).Msg(msg)
				respondWithGQLError(rw, ErrAuthentication)
				return
			}

			var identity primitives.Identity
			if typ == primitives.UserType {
				user := &primitives.User{}
				err = db.Get(ctx, session.IdentityID, user)
				if err != nil {
					msg := "Could not authenticate. Unable to find identity"
					logger.Error().Err(err).Msg(msg)
					respondWithGQLError(rw, ErrAuthentication)
					return
				}
				identity = user
			} else if typ == primitives.ServicePrimitiveType {
				service := &primitives.Service{}
				err = db.Get(ctx, session.IdentityID, service)
				if err != nil {
					msg := "Could not authenticate. Unable to find identity"
					logger.Error().Err(err).Msg(msg)
					respondWithGQLError(rw, ErrAuthentication)
					return
				}
				identity = service
			}

			ctx = context.WithValue(ctx, SessionContextKey, session)
			ctx = context.WithValue(ctx, IdentityContextKey, identity)

			next.ServeHTTP(rw, req)
		})
	}
}

func getGraphqlRoute(req *http.Request) (string, error) {
	params := &graphql.RawParams{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, params)
	if err != nil {
		return "", err
	}

	doc, gerr := parser.ParseQuery(&ast.Source{Input: params.Query})
	if gerr != nil {
		return "", err
	}

	return doc.Operations[0].SelectionSet[0].(*ast.Field).Name, nil
}

// RoundtripLoggerMiddleware logs information about request and response
// generated by the server. It depends on the logMiddleware.
func RoundtripLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, rw, req)
		logger := Logger(req.Context())

		logger.Info().
			Int("status", metrics.Code).
			Str("method", req.Method).
			Str("uri", req.RequestURI).
			Int64("duration", int64(metrics.Duration.Seconds()/1000)).
			Int64("size", metrics.Written).
			Str("referer", req.Referer()).
			Str("user-agent", req.UserAgent()).
			Msg("Finished")
	})
}

// RecoveryMiddleware catches any panics that occur in the call chain of the
// http request and response. If a panic does occur the panic is captured, a
// log is produced, and an internal server error is returned to the caller.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger := Logger(req.Context())
		defer func() {
			if p := recover(); p != nil {
				// If the panic contains an error then we want to log that
				// error otherwise we want to log the fact that we encountered
				// an error.
				//
				// XXX: It'd be helpful here to capture some portion of the
				// stack trace and log that as well.
				e, ok := p.(error)
				if !ok {
					e = errors.New(errors.UnknownCause, "Encountered panic: %s", p)
				}
				logger.Err(e).Msg("Encountered a panic; responding with a 500 error")

				// The following is the error we want to propagate externally
				err := errors.New(errors.UnknownCause, "Internal Server Error")
				respondWithError(rw, err)
			}
		}()

		next.ServeHTTP(rw, req)
	})
}

// respondWithError is a middleware helper for responding to an http request
// with a specific error. The error is written out in JSON and must be a
// partyerrors.Error.
func respondWithError(rw http.ResponseWriter, err error) {
	e := errors.ToError(err)
	respondWithJSON(rw, e.StatusCode(), e)
}

func respondWithGQLError(rw http.ResponseWriter, err error) {
	e := errors.ToError(err)

	// The context below is not being used the called function
	res := graphql.ErrorResponse(context.Background(), strings.Join(e.Messages, ","))
	respondWithJSON(rw, e.StatusCode(), res)
}

// respondWithJSON is a middleware helper for responding to an http request
// with a specific json response. If an error is encountered a log is produced.
func respondWithJSON(rw http.ResponseWriter, code int, out interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(out); err != nil {
		panic(fmt.Sprintf("Could not marshal json to response: %s", err))
	}
}
