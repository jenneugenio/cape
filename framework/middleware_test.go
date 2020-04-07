package framework

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/justinas/alice"
	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
	"github.com/rs/zerolog"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var logger *zerolog.Logger

func init() {
	logger = TestLogger()
}

func runMiddleware(mw http.Handler, header *http.Header) {
	req := httptest.NewRequest("GET", "http://my.capeprivacy.com", nil)
	if header != nil {
		req.Header = *header
	}
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)
}

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("sets the id on the request context", func(t *testing.T) {
		gm.RegisterTestingT(t)

		wasCalled := false
		next := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			wasCalled = true
			v := r.Context().Value(RequestIDContextKey)

			gm.Expect(v).ToNot(gm.BeNil())
			gm.Expect(rw.Header().Get("X-Request-Id")).To(gm.Equal(v.(uuid.UUID).String()))
		})

		mw := RequestIDMiddleware(next)
		runMiddleware(mw, nil)

		gm.Expect(wasCalled).To(gm.BeTrue(), "next was not called")
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	t.Run("returns a 500 error", func(t *testing.T) {
		gm.RegisterTestingT(t)

		wasCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			wasCalled = true
			panic("a bad thing happened")
		})

		mw := alice.New(RequestIDMiddleware, LogMiddleware(logger), RecoveryMiddleware).Then(next)

		req := httptest.NewRequest("GET", "http://my.cape.com", nil)
		w := httptest.NewRecorder()

		mw.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		gm.Expect(wasCalled).To(gm.BeTrue(), "next was not called")
		gm.Expect(resp.StatusCode).To(gm.Equal(http.StatusInternalServerError))
		gm.Expect(resp.Header.Get("Content-Type")).To(gm.Equal("application/json"))

		e := errors.Error{}
		err := json.Unmarshal(body, &e)
		gm.Expect(err).To(gm.BeNil(), "could not unmarshal into Error struct")
		gm.Expect(e.Cause).To(gm.Equal(errors.UnknownCause))
		gm.Expect(e.Messages[0]).To(gm.Equal("Internal Server Error"))
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("sets the logger on the request context", func(t *testing.T) {
		gm.RegisterTestingT(t)

		wasCalled := false
		next := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			wasCalled = true
			gm.Expect(func() {
				_ = Logger(r.Context())
			}).ToNot(gm.Panic())
		})

		mw := alice.New(RequestIDMiddleware, LogMiddleware(logger)).Then(next)
		runMiddleware(mw, nil)

		gm.Expect(wasCalled).To(gm.BeTrue(), "next was not called")
	})

	t.Run("panics if request id middleware is not before it", func(t *testing.T) {
		gm.RegisterTestingT(t)

		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
		mw := alice.New(LogMiddleware(logger)).Then(next)

		gm.Expect(func() {
			runMiddleware(mw, nil)
		}).To(gm.Panic())
	})
}

func TestAuthTokenMiddleware(t *testing.T) {
	t.Run("sets auth token on the request context", func(t *testing.T) {
		gm.RegisterTestingT(t)

		expID := base64.New([]byte("cool-auth-token"))
		wasCalled := false
		next := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			wasCalled = true
			id := r.Context().Value(AuthTokenContextKey)
			gm.Expect(id).To(gm.Equal(expID))
		})

		header := &http.Header{}
		header.Add("Authorization", "Bearer "+expID.String())

		mw := AuthTokenMiddleware(next)
		runMiddleware(mw, header)

		gm.Expect(wasCalled).To(gm.BeTrue(), "next was not called")
	})

	t.Run("bad auth token", func(t *testing.T) {
		gm.RegisterTestingT(t)

		expID := base64.New([]byte("cool-auth-token"))
		wasCalled := false
		next := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			wasCalled = true
		})

		header := http.Header{}
		header.Add("Authorization", expID.String())

		mw := AuthTokenMiddleware(next)
		req := httptest.NewRequest("GET", "http://my.capeprivacy.com", nil)
		req.Header = header

		w := httptest.NewRecorder()

		mw.ServeHTTP(w, req)

		gm.Expect(wasCalled).To(gm.BeFalse())

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		gResp := &graphql.Response{}
		err := json.Unmarshal(body, gResp)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(gResp.Errors)).To(gm.Equal(1))

		gm.Expect(gResp.Errors[0].Message).To(gm.Equal("Unable to parse auth header"))
	})
}
