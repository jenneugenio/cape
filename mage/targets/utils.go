package targets

import (
	"context"
	"os"

	"github.com/capeprivacy/cape/mage"
)

// We optionally support specifying the version via an env variable. This
// gets around us having to pull the version from Git which is useful when
// we're building Cape inside of a container without a checkout of git.
//
// If the environment variable is not set we then fall back to git!
func getVersion(ctx context.Context) (*mage.Version, error) {
	version := os.Getenv("VERSION")
	if len(version) > 0 {
		return mage.NewVersion(version)
	}

	deps, err := mage.Dependencies.Get([]string{"git"})
	if err != nil {
		return nil, err
	}

	git := deps[0].(*mage.Git)
	if err := git.Check(ctx); err != nil {
		return nil, err
	}

	return git.Tag(ctx)
}
