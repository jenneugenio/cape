package build

import (
	"context"
	"fmt"
	"sync"
)

// Dependencies contains a list of external dependencies
var Dependencies = deps{}

func init() {
	godep, err := NewGo()
	if err != nil {
		panic(err)
	}

	dockerdep, err := NewDocker()
	if err != nil {
		panic(err)
	}

	err = Dependencies.Add(godep)
	if err != nil {
		panic(err)
	}

	err = Dependencies.Add(dockerdep)
	if err != nil {
		panic(err)
	}
}

// RunnerFunc is a function that performs an actions on a dependency
type RunnerFunc func(d Dependency) error

// Dependency is an interface for an external dependency required for a build
// function. An example of a dependency is Go, Docker, or Protoc.
type Dependency interface {
	// Methods for checking if the required dependency
	Check(context.Context) error
	Name() string

	// Methods for setting up and tearing down anything needed or produced by
	// the dependency
	Setup(context.Context) error
	Clean(context.Context) error
}

// deps represents a registry of dependencies. It offers convenience
// functions for checking if a series of deps are met.
type deps map[string]Dependency

func (d deps) List() []string {
	out := make([]string, len(d))
	i := 0
	for k := range d {
		out[i] = k
		i++
	}

	return out
}

func (d deps) Get(names []string) ([]Dependency, error) {
	deps := make([]Dependency, len(names))
	for i, n := range names {
		dep, ok := d[n]
		if !ok {
			return deps, fmt.Errorf("Unknown dependency: %s", n)
		}

		deps[i] = dep
	}

	return deps, nil
}

// Run executes the given function against all of the named dependencies in parallel.
func (d deps) Run(names []string, f RunnerFunc) error {
	deps, err := d.Get(names)
	if err != nil {
		return err
	}

	errors := NewSafeSlice()
	wg := &sync.WaitGroup{}

	for _, dep := range deps {
		item := dep
		wg.Add(1)
		go func() {
			defer wg.Done()
			errors.Add(f(item))
		}()
	}

	wg.Wait()
	// TODO: Introduce a "multi-error" type enabling us to return multiple
	// errors, cleanly.
	for _, err := range errors.Get() {
		if err != nil {
			return err.(error)
		}
	}

	return nil
}

func (d deps) Add(dep Dependency) error {
	if _, ok := d[dep.Name()]; ok {
		return fmt.Errorf("Cannot add dependency, it already exists: %s", dep.Name())
	}

	d[dep.Name()] = dep
	return nil
}
