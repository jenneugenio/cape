package targets

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"

	"github.com/capeprivacy/cape/mage"
)

var cluster = &mage.Cluster{
	Name:       "cape-local",
	ConfigFile: "build/local/kind/cape-local.yaml",
}

var registry = &mage.Registry{
	Name: "cape-local-docker-registry",
	Port: "5000",
}

func init() {
	err := mage.Tracker.Add(fmt.Sprintf("kind://%s", cluster.Name), mage.CleanKind(cluster))
	if err != nil {
		panic(err)
	}

	err = mage.Tracker.Add(fmt.Sprintf("registry://%s", registry.Name), mage.CleanDockerRegistry(registry))
	if err != nil {
		panic(err)
	}
}

type Local mg.Namespace

// Create creates a local kubernetes cluster, builds the required docker
// images, and then deploys their subsequent helm packages into the cluster.
func (l Local) Create(ctx context.Context) error {
	required := []string{"kind", "docker_registry"}
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
	dockerRegistry := deps[1].(*mage.DockerRegistry)

	network, err := dockerRegistry.Create(ctx, registry)
	if err != nil {
		return err
	}

	err = kind.Create(ctx, cluster, registry)
	if err != nil {
		return err
	}

	if _, ok := network.Networks["kind"]; ok {
		return nil
	}

	// Need to connecto the docker registry to the kind docker network
	return dockerRegistry.Connect(ctx, registry, "kind")
}

// Deploy builds and deploys cape from your local repository to the local
// kubernetes cluster. If a cluster is not running one will be created.
func (l Local) Deploy(ctx context.Context) error {
	mg.Deps(Local.Create, Build.Docker)
	return nil
}

// Status returns the current status of the kubernetes cluster, services, and
// jobs that are deployed by cape into the local cape cluster.
func (l Local) Status(ctx context.Context) error {
	required := []string{"kind", "docker_registry"}
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
	dockerRegistry := deps[1].(*mage.DockerRegistry)

	kindState, err := kind.Status(ctx, cluster)
	if err != nil {
		return err
	}

	registryState, err := dockerRegistry.Status(ctx, registry)
	if err != nil {
		return err
	}

	fmt.Printf("Kind Cluster:\t%s\n", kindState)
	fmt.Printf("Docker Registry:\t%s\n", registryState)

	return nil
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
