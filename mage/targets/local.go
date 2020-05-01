package targets

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"

	"github.com/capeprivacy/cape/mage"
)

var cluster = &mage.Cluster{
	Name:       "cape-local",
	ConfigFile: "mage/config/kind/cape-local.yaml",
}

func init() {
	err := mage.Tracker.Add(fmt.Sprintf("kind://%s", cluster.Name), mage.CleanKind(cluster))
	if err != nil {
		panic(err)
	}
}

type Local mg.Namespace

// Create creates a local kubernetes cluster, builds the required docker
// images, and then deploys their subsequent helm packages into the cluster.
func (l Local) Create(ctx context.Context) error {
	required := []string{"kind"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	deps, err := mage.Dependencies.Get(required)
	if err != nil {
		return err
	}

	kind := deps[0].(*mage.Kind)
	return kind.Create(ctx, cluster)
}

// Status returns the current status of the kubernetes cluster, services, and
// jobs that are deployed by cape into the local cape cluster.
func (l Local) Status(ctx context.Context) error {
	required := []string{"kind"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	deps, err := mage.Dependencies.Get(required)
	if err != nil {
		return err
	}

	kind := deps[0].(*mage.Kind)
	return kind.Status(ctx, cluster)
}

// Destroy deletes the kubernetes clusters and any managed volumes completely
// erasing anything related to the local deployment
func (l Local) Destroy(ctx context.Context) error {
	required := []string{"kind"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	deps, err := mage.Dependencies.Get(required)
	if err != nil {
		return err
	}

	kind := deps[0].(*mage.Kind)
	return kind.Destroy(ctx, cluster)
}
