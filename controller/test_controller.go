package controller

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/dbtest"
	"github.com/dropoutlabs/cape/framework"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

var (
	TimeoutCause = errors.NewCause(errors.InternalServerErrorCategory, "server_start_timeout")
)

// TestController is a convenience wrapper around the controller to help with testing.
// In integration testing, we often need to spin up a database, migrate it, run the tests, tear down, and then stop
// the controller.  This amounts to a lot of boilerplate at the top of tests.  You can wrap all of that boilerplate up
// with this simple struct
//
//     // handle your errors!!
//     tc, err := controller.NewTestController()
//     _, err = tc.Setup(ctx) // this can return an actual controller if you need one
//     defer tc.Teardown(ctx)
type TestController struct {
	controller *Controller
	database   dbtest.TestDatabase
	logger     *zerolog.Logger

	User *primitives.User

	// TODO this should be stored on our own client
	Credentials *auth.Credentials
}

// NewTestController gives you a test controller with a live database & gql server
func NewTestController() (*TestController, error) {
	testDB, err := dbtest.New(os.Getenv("CAPE_DB_URL"))
	if err != nil {
		return nil, err
	}

	logger := framework.TestLogger()
	controller, err := New(testDB.URL(), logger, "cape-test-controller")
	if err != nil {
		return nil, err
	}

	return &TestController{
		controller: controller,
		database:   testDB,
		logger:     logger,
	}, nil
}

// Setup will run db migrations, connect to the db, and start the gql server
func (t *TestController) Setup(ctx context.Context) (*Controller, error) {
	err := t.database.Setup(ctx)
	if err != nil {
		return nil, err
	}

	migrations := []string{
		os.Getenv("CAPE_DB_MIGRATIONS"),
		os.Getenv("CAPE_DB_TEST_MIGRATIONS"),
	}

	migrator, err := database.NewMigrator(t.database.URL(), migrations...)
	if err != nil {
		return nil, err
	}

	err = migrator.Up(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		err := t.controller.Start(ctx)
		if err != nil {
			t.logger.Error().Err(err).Msg("Could not start test controller")
		}
	}()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()

	for {
		// We are never bubbling this error up to the caller, but that is intentional
		// This request will fail until the server is online, we will ping /health every 50ms until we get a 200
		// If 5s elapses then we will give up and fail whatever test is using this.
		resp, err := http.Get("http://localhost:8081/_healthz")
		if err == nil {
			if resp.StatusCode == 200 {
				err = t.createUser(ctx)
				if err != nil {
					return nil, err
				}

				return t.controller, nil
			}
		}

		select {
		case <-timeout:
			return nil, errors.New(TimeoutCause, "Timed out waiting for the server to start")
		case <-time.After(50 * time.Millisecond):
			continue
		}
	}
}

// Teardown will destroy the database and stop the gql server
func (t *TestController) Teardown(ctx context.Context) error {
	err := t.controller.Stop(ctx)
	if err != nil {
		return err
	}

	return t.database.Teardown(ctx)
}

func (t *TestController) createUser(ctx context.Context) error {
	creds, err := auth.NewCredentials([]byte("swimmingfishmustswimquick"), nil)
	if err != nil {
		return err
	}
	t.Credentials = creds

	user, err := primitives.NewUser("Testy Testerson", "testy@capeprivacy.com", creds.Package())
	if err != nil {
		return err
	}

	t.User = user

	return t.controller.backend.Create(ctx, user)
}
