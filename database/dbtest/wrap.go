package dbtest

import (
	"context"
	"sync"
)

// Wrapper is a TestDatabase providing a singleton-style experience allowing
// more than one test to use the same underlying database. The purpose is to
// save time on how long it takes to setup a database between tests.
//
// Wrappers are kept in a centralized managed by the New function.
type Wrapper struct {
	db    TestDatabase
	mutex *sync.Mutex
	count int
}

// Wrap returns a wrapped test database setting it up to be a singleton-style
// object.
func Wrap(db TestDatabase) TestDatabase {
	return &Wrapper{
		db:    db,
		mutex: &sync.Mutex{},
		count: 0,
	}
}

// Setup calls the setup method on the TestDatabase if this is the first
// invocation
func (w *Wrapper) Setup(ctx context.Context) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	err := w.db.Setup(ctx)
	if err != nil {
		return err
	}

	w.count += 1
	return nil
}

// Teardown calls the teardown method on the TestDatabase if this is the last
// invocation of the method.
func (w *Wrapper) Teardown(ctx context.Context) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.count > 2 {
		w.count -= 1
		return nil
	}

	err := w.db.Teardown(ctx)
	w.count = 0
	return err
}

// Truncate is just a pass through to the TestDatabase's truncation method. It
// does not prevent truncation from being called more than once.
//
// It's assumed this responsibility is managed by the underlying database.
func (w *Wrapper) Truncate(ctx context.Context) error {
	return w.db.Truncate(ctx)
}

// URL returns the underlying URL of the Test Database
func (w *Wrapper) URL() string {
	return w.db.URL()
}

// Database returns a reference to the underlying TestDatabase
func (w *Wrapper) Database() TestDatabase {
	return w.db
}
