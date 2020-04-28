package mage

import (
	"sync"

	"go.uber.org/multierr"
)

// Errors is a concurrency safe wrapper around uber's multierr
type Errors struct {
	errors []error
	lock   *sync.Mutex
}

func NewErrors() *Errors {
	return &Errors{
		errors: []error{},
		lock:   &sync.Mutex{},
	}
}

// Error completes the error interface
func (e *Errors) Error() string {
	e.lock.Lock()
	defer e.lock.Unlock()

	return multierr.Combine(e.errors...).Error()
}

// Err returns the underlying error (which can be nil if no non-nil errors were
// appended)
func (e *Errors) Err() error {
	return multierr.Combine(e.errors...)
}

// Errors returns the underlying slice of errors
func (e *Errors) Errors() []error {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.errors
}

// Append adds the error into the underlying multierr combining all of the
// errors together into one error
func (e *Errors) Append(err error) bool {
	if err == nil {
		return false
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	e.errors = append(e.errors, err)
	return true
}
