package targets

import (
	"context"

	"github.com/capeprivacy/cape/mage"
)

// Bootstrap installs any required tooling/setup enabling a user to build and
// deploy Cape locally for development and testing.
func Bootstrap(ctx context.Context) error {
	// TODO: Look into making a UI or similar (e.g. printing success output)
	// TODO: Add warning that this could take time on the first run.
	deps := []string{"go"} // TODO: Make Type Safe & Compile Time
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

// Check checks whether or not the environment contains all of the required
// dependencies to bootstrap, test, build, and deploy Cape locally.
func Check(ctx context.Context) error {
	return mage.Dependencies.Run(mage.Dependencies.List(), func(d mage.Dependency) error {
		return d.Check(ctx)
	})
}

// Clean removes any installed tools, modules, or build artifacts created by
// any targets.
//
// This command will remove 'Magefile' which will need to be installed again if
// you run it multiple times via the bootstrap.go command.
func Clean(ctx context.Context) error {
	deps := mage.Dependencies.List()

	// We collect all of the errors that occurred and don't report them until
	// the end. This way we can clean up everything to the best of our ability.
	errors := mage.NewErrors()
	errors.Append(mage.Tracker.Clean(ctx))
	errors.Append(mage.Dependencies.Run(deps, func(d mage.Dependency) error {
		return d.Clean(ctx)
	}))

	return errors.Err()
}
