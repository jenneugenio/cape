package mage

import (
	"context"
	"fmt"
	"sync"
)

// Dependencies contains a list of external dependencies
var Dependencies = deps{}

// TODO: Change this from a global init function to something thats created for
// each magefile command invocation. That way, it's less hard coding, easier to
// test!
func init() {
	dockerdep := MustDocker("18.0")

	depVersions := []Dependency{
		MustGit("2.0"),
		dockerdep,
		MustKind(dockerdep, MustKubectl("1.11"), "0.8"),
		MustDockerRegistry(dockerdep),
		MustGolang("github.com/capeprivacy/cape", "1.14"),
		MustHelm("3.2", map[string]string{
			"bitnami": "https://charts.bitnami.com/bitnami",
		}),
	}

	for _, d := range depVersions {
		Dependencies.MustAdd(d)
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

	errors := NewErrors()
	wg := &sync.WaitGroup{}

	for _, dep := range deps {
		item := dep
		wg.Add(1)
		go func() {
			defer wg.Done()
			errors.Append(f(item))
		}()
	}

	wg.Wait()

	return errors.Err()
}

func (d deps) Add(dep Dependency) error {
	if _, ok := d[dep.Name()]; ok {
		return fmt.Errorf("Cannot add dependency, it already exists: %s", dep.Name())
	}

	d[dep.Name()] = dep
	return nil
}

func (d deps) MustAdd(dep Dependency) {
	if err := d.Add(dep); err != nil {
		panic(err)
	}
}
