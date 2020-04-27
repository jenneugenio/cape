// +build mage
package base

import (
	"context"

	"github.com/capeprivacy/cape/mage"
)

// Bootstrap installs any required tooling/setup enabling a user to build and
// deploy Cape locally for development and testing.
func Bootstrap(ctx context.Context) error {
	// TODO: Look into making a UI or similar (e.g. printing success output)
	// TODO: Add warning that this could take time on the first run.
	deps := []string{"go", "docker"} // TODO: Make Type Safe & Compile Time
	err := mage.Dependencies.Run(deps, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	return mage.Dependencies.Run(deps, func(d mage.Dependency) error {
		return d.Setup(ctx)
	})
}

// Clean removes any installed tools, modules, or build artifacts created by
// any targers.
//
// This command will remove 'Magefile' which will need to be installed again if
// you run it multiple times via the bootstrap.go command.
func Clean(ctx context.Context) error {
	deps := mage.Dependencies.List()
	return mage.Dependencies.Run(deps, func(d mage.Dependency) error {
		return d.Clean(ctx)
	})
}
