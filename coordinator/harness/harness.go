package harness

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/database/dbtest"
	"github.com/capeprivacy/cape/framework"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

var (
	TimeoutCause    = errors.NewCause(errors.RequestTimeoutCategory, "start_timeout")
	NotStartedCause = errors.NewCause(errors.BadRequestCategory, "coordinator_not_started")
)

// Harness represents a http server used for testing. Its responsibility is to
// provide the transport layer for the application contained within the
// component. In this case, the component represents a Coordinator.
//
// This is a convenience layer for testing the Coordinator by spinning up a
// database, migrating it, running the tests, tearing them down, and then
// cleaning any remaining artifacts up.
//
// This harness uses net/http/httptest.Server so it will start the server on a
// randomized port that has not yet been assigned.
//
// This harness _does not_ manage any client state or actually "setup" the
// coordinator admin functionality. Please see the harness.Manager which
// provides convenience functions for managing application state.
//
// You can use it as follows:
//
//  h, err := NewHarness(cfg)
//  if err != nil {
//    return err // Handle your errors!
//  }
//  err = h.Setup(ctx)
//  ...
//  defer h.Teardown(ctx)
//
type Harness struct {
	db        dbtest.TestDatabase
	server    *httptest.Server
	component framework.Component
	logger    *zerolog.Logger
	manager   *Manager
	cfg       *Config
}

// NewHarness returns a new harness that's configured and ready to be setup!
func NewHarness(cfg *Config) (*Harness, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Harness{
		cfg: cfg,
	}, nil
}

// Start sets up the testing harness to test the Coordinator component
func (h *Harness) Setup(ctx context.Context) error {
	logger := framework.TestLogger()

	db, err := dbtest.New(h.cfg.dbURL)
	if err != nil {
		return err
	}

	err = db.Setup(ctx)
	if err != nil {
		return err
	}

	cleanupWasCalled := false
	cleanup := func(in error) error {
		// XXX: Should we return a different error?
		if cleanupWasCalled {
			return in
		}

		cleanupWasCalled = true
		if h.server != nil {
			h.server.Close() //httptest.Server.Close() does not return an error
		}

		if h.component != nil {
			err := h.component.Teardown(ctx)
			logger.Error().Msgf("Could not stop coordinator component: %s", err)
		}

		h.db = nil
		h.server = nil
		h.component = nil
		h.manager = nil

		err := db.Teardown(ctx)
		if err != nil {
			return err
		}

		return in
	}

	migrator, err := database.NewMigrator(db.URL(), h.cfg.Migrations()...)
	if err != nil {
		return cleanup(err)
	}

	err = migrator.Up(ctx)
	if err != nil {
		return cleanup(err)
	}

	dbURL, err := primitives.DBURLFromURL(db.URL())
	if err != nil {
		return cleanup(err)
	}

	kp, err := auth.NewKeypair()
	if err != nil {
		return cleanup(err)
	}

	coordinator, err := coordinator.New(&coordinator.Config{
		DB: &coordinator.DBConfig{
			Addr: dbURL,
		},
		InstanceID: "cape",
		Port:       1, // This port is ignored!
		Auth: &coordinator.AuthConfig{
			KeypairPackage: kp.Package(),
		},
	}, logger)
	if err != nil {
		return err
	}

	handler, err := coordinator.Setup(ctx)
	if err != nil {
		return cleanup(err)
	}

	h.logger = logger
	h.component = coordinator
	h.db = db

	// httptest.NewServer starts listening immediately, it also picks a
	// randomized port to listen on!
	h.server = httptest.NewServer(handler)
	h.manager = &Manager{h: h}

	// We try to wait for the coordinator to start for _up to_ 5 seconds! At
	// which point we bail out and return an error.
	timeout := time.NewTimer(5 * time.Second)
	for {
		// We are never bubbling this error up to the caller, but that is
		// intentional This request will fail until the server is online, we
		// will ping /_healthz every 50ms until we get a 200 If 5s elapses then
		// we will give up and fail whatever test is using this.
		u, err := h.URL()
		if err != nil {
			return cleanup(err)
		}
		resp, err := http.Get(u.String() + "/_healthz")
		if err == nil {
			if resp.StatusCode == 200 {
				if err != nil {
					return cleanup(err)
				}

				return nil
			}
		}

		select {
		case <-timeout.C:
			return cleanup(errors.New(TimeoutCause, "Timed out waiting for coordinator to start"))
		case <-time.After(50 * time.Millisecond):
			continue
		}
	}
}

// Teardown destroys all of the underlying resources needed by the Coordinator
// and stops the test server from serving it at all!
func (h *Harness) Teardown(ctx context.Context) error {
	if h.component == nil || h.db == nil || h.server == nil {
		return errors.New(NotStartedCause, "Harness must be started to be torn down")
	}

	h.server.Close()

	err := h.component.Teardown(ctx)
	if err != nil {
		h.logger.Error().Msgf("Could not cleanly stop coordinator component: %s", err)
	}

	db := h.db
	h.db = nil
	h.component = nil
	h.server = nil
	h.manager = nil

	err = db.Teardown(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Client returns an unauthenticated Client for the underlying instance of the
// coordinator.
func (h *Harness) Client() (*coordinator.Client, error) {
	u, err := h.URL()
	if err != nil {
		return nil, err
	}

	transport := coordinator.NewTransport(u, nil)
	return coordinator.NewClient(transport), nil
}

// Manager returns a test state manager for this Harness
func (h *Harness) Manager() *Manager {
	return h.manager
}

// URL returns the url to the running coordinator once the harness has been started.
func (h *Harness) URL() (*primitives.URL, error) {
	if h.server == nil {
		return nil, errors.New(NotStartedCause, "Harness must be started to retrieve url")
	}

	return primitives.NewURL(h.server.URL)
}
