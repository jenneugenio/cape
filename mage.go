// +build mage

package main

import (
	"context"

	"github.com/capeprivacy/cape/build"
)

// Bootstrap installs any required tooling/setup enabling a user to build and
// deploy Cape locally for development and testing.
func Bootstrap(ctx context.Context) error {
	// TODO: Look into making a UI or similar (e.g. printing success output)
	// TODO: Add warning that this could take time on the first run.
	deps := []string{"go", "docker"} // TODO: Make Type Safe & Compile Time
	err := build.Dependencies.Run(deps, func(d build.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	return build.Dependencies.Run(deps, func(d build.Dependency) error {
		return d.Setup(ctx)
	})
}

// Clean removes any installed tools, modules, or build artifacts created by
// any targers.
//
// This command will remove 'Magefile' which will need to be installed again if
// you run it multiple times via the bootstrap.go command.
func Clean(ctx context.Context) error {
	deps := build.Dependencies.List()
	return build.Dependencies.Run(deps, func(d build.Dependency) error {
		return d.Clean(ctx)
	})
}
