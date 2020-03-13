package controller

import (
	"context"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/dbtest"
	"os"
	"time"
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
}

// NewTestController gives you a test controller with a live database & gql server
func NewTestController() (*TestController, error) {
	testDB, err := dbtest.New(os.Getenv("CAPE_DB_URL"))
	if err != nil {
		return nil, err
	}

	controller, err := New(testDB.URL(), "test-controller")
	if err != nil {
		return nil, err
	}

	return &TestController{
		controller: controller,
		database:   testDB,
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

	t.controller.Start(ctx)

	// TODO -- Delete me
	// This should be removed once we have health checks, ping the health endpoint until we get a 200
	time.Sleep(2 * time.Second)
	return t.controller, nil
}

// Teardown will destroy the database and stop the gql server
func (t *TestController) Teardown(ctx context.Context) error {
	err := t.controller.Stop(ctx)
	if err != nil {
		return err
	}

	return t.database.Teardown(ctx)
}
