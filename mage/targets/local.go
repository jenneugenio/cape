package targets

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

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

var charts = []*mage.Chart{
	{
		Name:    "postgres",
		Chart:   "bitnami/postgresql",
		Version: "8.9.4",
		Values:  "mage/config/postgres-values.yaml",
	},
	{
		Name:    "coordinator",
		Chart:   "charts/coordinator",
		Version: "0.0.1",
		Values:  "mage/config/coordinator-values.yaml",
	},
	{
		Name:    "connector",
		Chart:   "charts/connector",
		Version: "0.0.1",
		Values:  "mage/config/connector-values.yaml",
	},
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
	required := []string{"kind", "docker_registry", "helm"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	err = mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Setup(ctx)
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

	// Need to connect to the docker registry to the kind docker network
	return dockerRegistry.Connect(ctx, registry, "kind")
}

func (l Local) Push(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, Build.Docker)

	for _, image := range dockerImages {
		newTag := "localhost:5000/" + image.String()
		tagCmd := []string{"tag", image.String(), newTag}
		err := sh.RunV("docker", tagCmd...)
		if err != nil {
			return err
		}

		pushCmd := []string{"push", newTag}
		err = sh.RunV("docker", pushCmd...)
		if err != nil {
			return err
		}
	}

	return nil
}

// Deploy builds and deploys cape from your local repository to the local
// kubernetes cluster. If a cluster is not running one will be created.
func (l Local) Deploy(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, Local.Create, Local.Push)

	required := []string{"helm"}
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

	helm := deps[0].(*mage.Helm)
	errors := mage.NewErrors()

	for _, chart := range charts {
		errors.Append(helm.Install(ctx, chart))
	}

	return errors.Err()
}

// Status returns the current status of the kubernetes cluster, services, and
// jobs that are deployed by cape into the local cape cluster.
func (l Local) Status(ctx context.Context) error {
	required := []string{"kind", "docker_registry", "helm"}
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
	helm := deps[2].(*mage.Helm)

	kindState, err := kind.Status(ctx, cluster)
	if err != nil {
		return err
	}

	registryState, err := dockerRegistry.Status(ctx, registry)
	if err != nil {
		return err
	}

	releases, err := helm.List(ctx)
	if err != nil {
		return err
	}

	releaseMap := map[string]*mage.Release{}
	for _, release := range releases {
		releaseMap[release.Name] = release
	}

	fmt.Printf("Kind:\t\t%s\n", kindState)
	fmt.Printf("Registry:\t%s\n", registryState)

	fmt.Printf("\nHelm Charts:\n")
	for _, chart := range charts {
		release, ok := releaseMap[chart.Name]
		status := "unknown"
		if ok {
			status = release.Status
		}

		fmt.Printf("\t%s: %s\n", chart.Name, status)
	}

	fmt.Printf("\nRun `kubectl get svc -o wide` to check on the status of the underlying pods for the deployed charts.\n")
	return nil
}

// Destroy deletes the kubernetes clusters and any managed volumes completely
// erasing anything related to the local deployment
func (l Local) Destroy(ctx context.Context) error {
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

	errors := mage.NewErrors()
	errors.Append(kind.Destroy(ctx, cluster))
	errors.Append(dockerRegistry.Destroy(ctx, registry))
	return errors.Err()
}
