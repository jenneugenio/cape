package harness

import (
	"context"
	"crypto/x509"
	"net/http/httptest"
	"time"

	"github.com/manifoldco/go-base64"
	"github.com/rs/zerolog"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/connector"
	"github.com/dropoutlabs/cape/connector/client"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/dbtest"
	"github.com/dropoutlabs/cape/framework"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// ConnectorEmail is the email of the connector
const ConnectorEmail = "service:data-connector@cape.com"

var (
	// TimeoutCause the connector took to long to start
	TimeoutCause = errors.NewCause(errors.RequestTimeoutCategory, "start_timeout")

	// NotStartedCause the connector has not started yet
	NotStartedCause = errors.NewCause(errors.BadRequestCategory, "connector_not_started")
)

// Harness represents a http server used for testing. Its responsibility is to
// provide the transport layer for the application contained within the
// component. In this case, the component represents a connector.
//
// This is a convenience layer for testing the connector by spinning up a
// database, migrating it, running the tests, tearing them down, and then
// cleaning any remaining artifacts up.
//
// This harness uses net/http/httptest.Server so it will start the server on a
// randomized port that has not yet been assigned.
//
// You can use it as follows:
//
//  h, err := NewHarness()
//  if err != nil {
//    return err // Handle your errors!
//  }
//  err = h.Setup(ctx)
//  ...
//  defer h.Teardown(ctx)
//
type Harness struct {
	cfg       *Config
	server    *httptest.Server
	db        dbtest.TestDatabase
	component framework.Component
	logger    *zerolog.Logger
	apiToken  *auth.APIToken
}

// NewHarness returns a new harness that's configured and ready to be setup!
func NewHarness(cfg *Config) (*Harness, error) {
	email, err := primitives.NewEmail(ConnectorEmail)
	if err != nil {
		return nil, err
	}

	apiToken, err := auth.NewAPIToken(email, cfg.ControllerURL)
	if err != nil {
		return nil, err
	}

	return &Harness{
		cfg:      cfg,
		apiToken: apiToken,
	}, nil
}

// Setup sets up the testing harness to test the connector component
func (h *Harness) Setup(ctx context.Context) error {
	logger := framework.TestLogger()

	db, err := dbtest.NewTestPostgres(h.cfg.dbURL.String())
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
			logger.Error().Msgf("Could not stop connector component: %s", err)
		}

		h.server = nil
		h.component = nil

		err := db.Teardown(ctx)
		if err != nil {
			return err
		}

		return in
	}

	migrator, err := database.NewMigrator(db.URL(), h.cfg.sourceMigrationsDir)
	if err != nil {
		return cleanup(err)
	}

	err = migrator.Up(ctx)
	if err != nil {
		return cleanup(err)
	}

	connector, err := connector.New(&connector.Config{
		InstanceID: "cape-connector",
		Port:       1, // This port is ignored!
		Token:      h.apiToken,
	}, logger)
	if err != nil {
		return err
	}

	handler, err := connector.Setup(ctx)
	if err != nil {
		return cleanup(err)
	}

	h.logger = logger
	h.component = connector
	h.db = db

	// httptest.NewServer starts listening immediately, it also picks a
	// randomized port to listen on!
	h.server = httptest.NewUnstartedServer(handler)
	h.server.EnableHTTP2 = true
	h.server.StartTLS()

	client := h.server.Client()

	// We try to wait for the connector to start for _up to_ 5 seconds! At
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
		resp, err := client.Get(u.String() + "/_healthz")
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
			return cleanup(errors.New(TimeoutCause, "Timed out waiting for connector to start"))
		case <-time.After(50 * time.Millisecond):
			continue
		}
	}
}

// Teardown destroys all of the underlying resources needed by the connector
// and stops the test server from serving it at all!
func (h *Harness) Teardown(ctx context.Context) error {
	if h.component == nil || h.server == nil {
		return errors.New(NotStartedCause, "Harness must be started to be torn down")
	}

	h.server.Close()

	err := h.component.Teardown(ctx)
	if err != nil {
		h.logger.Error().Msgf("Could not cleanly stop connector component: %s", err)
	}

	db := h.db
	h.db = nil
	h.component = nil
	h.server = nil

	err = db.Teardown(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Client returns an unauthenticated Client for the underlying instance of the
// connector.
func (h *Harness) Client(authToken *base64.Value) (*client.Client, error) {
	u, err := h.URL()
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(h.server.Certificate())

	return client.NewClient(u, authToken, certPool)
}

// URL returns the url to the running connector once the harness has been started.
func (h *Harness) URL() (*primitives.URL, error) {
	if h.server == nil {
		return nil, errors.New(NotStartedCause, "Harness must be started to retrieve url")
	}

	return primitives.NewURL(h.server.URL)
}

// APIToken returns the APIToken needed by the connector
func (h *Harness) APIToken() *auth.APIToken {
	return h.apiToken
}

// SourceCredentials manages the source credentials
func (h *Harness) SourceCredentials() *primitives.DBURL {
	return &primitives.DBURL{URL: h.db.URL()}
}
