package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	gm "github.com/onsi/gomega"
	"github.com/rs/zerolog"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var logger *zerolog.Logger

func init() {
	logger = TestLogger()
}

func runMiddleware(mw http.Handler) {
	req := httptest.NewRequest("GET", "http://my.capeprivacy.com", nil)
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

		mw := requestIDMiddleware(next)
		runMiddleware(mw)

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
		mw := requestIDMiddleware(logMiddleware(logger, recoveryMiddleware(next)))

		req := httptest.NewRequest("GET", "http://api.torus.sh", nil)
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

		mw := requestIDMiddleware(logMiddleware(logger, next))
		runMiddleware(mw)

		gm.Expect(wasCalled).To(gm.BeTrue(), "next was not called")
	})

	t.Run("panics if request id middleware is not before it", func(t *testing.T) {
		gm.RegisterTestingT(t)

		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
		mw := logMiddleware(logger, next)

		gm.Expect(func() {
			runMiddleware(mw)
		}).To(gm.Panic())
	})
}
