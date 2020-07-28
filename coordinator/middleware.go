package coordinator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/felixge/httpsnoop"
	"github.com/gofrs/uuid"
	"github.com/manifoldco/go-base64"
	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/db"
	fw "github.com/capeprivacy/cape/framework"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// RequestIDMiddleware sets a UUID on the response header and request context
// for use in tracing and
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		id, err := uuid.NewV4()
		if err != nil {
			panic(fmt.Sprintf("Could not generate a v4 uuid: %s", err))
		}

		ctx := context.WithValue(req.Context(), fw.RequestIDContextKey, id)
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
			requestID := fw.RequestID(ctx)
			logger := log.With().Str("request_id", requestID.String()).Logger()

			ctx = context.WithValue(ctx, fw.LoggerContextKey, logger)
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

		cookie, err := req.Cookie("token")
		if err != nil {
			next.ServeHTTP(rw, req)
			return
		}

		token, err := base64.NewFromString(cookie.Value)
		if err != nil {
			respondWithError(rw, req.URL.Path, err)
			return
		}

		ctx = context.WithValue(ctx, fw.AuthTokenContextKey, token)
		req = req.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}

// RoundtripLoggerMiddleware logs information about request and response
// generated by the server. It depends on the logMiddleware.
func RoundtripLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, rw, req)
		logger := fw.Logger(req.Context())

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
		logger := fw.Logger(req.Context())
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

				logger.Trace().Msg(string(debug.Stack()))

				// The following is the error we want to propagate externally
				err := errors.New(errors.UnknownCause, "Internal Server Error")
				respondWithError(rw, req.URL.Path, err)
			}
		}()

		next.ServeHTTP(rw, req)
	})
}

// respondWithError is a middleware helper for responding to an http request
// with a specific error. The error is written out in JSON and must be a
// partyerrors.Error.
func respondWithError(rw http.ResponseWriter, path string, err error) {
	if path == "/v1/query" {
		respondWithGQLError(rw, err)
		return
	}
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

// IsAuthenticatedMiddleware checks to make sure a query is authenticated
func IsAuthenticatedMiddleware(coordinator *Coordinator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			ta := coordinator.tokenAuth
			db := coordinator.backend
			capedb := coordinator.db

			logger := fw.Logger(ctx)
			token := fw.AuthToken(ctx)

			if token == nil {
				logger.Info().Msg("Could not authenticate. Token missing")
				respondWithError(rw, req.URL.Path, auth.ErrAuthentication)
				return
			}

			id, err := ta.Verify(token)
			if err != nil {
				msg := "Could not authenticate. Unable to verify auth token"
				logger.Info().Err(err).Msg(msg)
				respondWithError(rw, req.URL.Path, auth.ErrAuthentication)
				return
			}

			session := &primitives.Session{}
			err = db.Get(ctx, id, session)
			if err != nil {
				msg := "Could not authenticate. Unable to find session"
				logger.Info().Err(err).Msg(msg)
				respondWithError(rw, req.URL.Path, auth.ErrAuthentication)
				return
			}

			cp, err := getCredentialsProvider(ctx, db, capedb, session.UserID)
			if err != nil {
				msg := "Could not authenticate. Unable get credentialProvider type"
				logger.Info().Err(err).Msg(msg)
				respondWithError(rw, req.URL.Path, auth.ErrAuthentication)
				return
			}

			user, err := capedb.Users().GetByID(ctx, cp.GetUserID())
			if err != nil {
				respondWithError(rw, req.URL.Path, err)
				return
			}

			roles, err := capedb.Roles().GetAll(ctx, cp.GetUserID())
			if err != nil {
				respondWithError(rw, req.URL.Path, err)
				return
			}

			aSession, err := auth.NewSession(user, session, *roles, cp)
			if err != nil {
				respondWithError(rw, req.URL.Path, err)
				return
			}

			logger = logger.With().Str("user_id", aSession.GetID()).Logger()

			ctx = context.WithValue(ctx, fw.LoggerContextKey, logger)
			ctx = context.WithValue(ctx, fw.SessionContextKey, aSession)

			req = req.WithContext(ctx)

			next.ServeHTTP(rw, req)
		})
	}
}

func getCredentialsProvider(ctx context.Context, db database.Backend, capedb db.Interface, id string) (primitives.CredentialProvider, error) {
	dID, err := database.DecodeFromString(id)
	if err != nil {
		user, err := capedb.Users().GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	token := &primitives.Token{}
	err = db.Get(ctx, dID, token)

	return token, err
}
