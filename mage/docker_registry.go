package mage

import (
	"context"
	"fmt"
	"time"

	"github.com/magefile/mage/sh"
)

// Registry represents a docker registry managed and created by the
// DockerRegistry struct.
type Registry struct {
	Name string
	Port string
}

// Docker Registry represents a dependency for creating a Docker Registry using
// a Docker container.
//
// Docker must be installed on the system.
type DockerRegistry struct {
	docker *Docker
}

func NewDockerRegistry(docker *Docker) (*DockerRegistry, error) {
	return &DockerRegistry{docker: docker}, nil
}

func MustDockerRegistry(docker *Docker) *DockerRegistry {
	return &DockerRegistry{docker: docker}
}

func (d *DockerRegistry) Name() string {
	return "docker_registry"
}

func (d *DockerRegistry) Check(ctx context.Context) error {
	return d.docker.Check(ctx)
}

func (d *DockerRegistry) Setup(_ context.Context) error {
	return nil
}

func (d *DockerRegistry) Clean(_ context.Context) error {
	return nil
}

// Status returns whether or not the container is currently running
func (d *DockerRegistry) Status(ctx context.Context, r *Registry) (Status, error) {
	return d.docker.Status(ctx, r.Name)
}

// Create starts a docker registry with the given name and port combination.
// The ip address of the created registry is returned
//
// If the given registry is already running it returns an error.
func (d *DockerRegistry) Create(ctx context.Context, r *Registry) (*NetworkSettings, error) {
	status, err := d.docker.Status(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	if status == Running {
		return d.docker.Network(ctx, r.Name)
	}

	port := fmt.Sprintf("%s:5000", r.Port)
	err = sh.Run("docker", "run", "-d", "--restart=always", "-p", port, "--name", r.Name, "registry:2")
	if err != nil {
		return nil, err
	}

	// Wait for the docker container to become "running" before attempting to
	// get it's address
	err = WaitFor(ctx, func(ctx context.Context) (bool, error) {
		status, err := d.Status(ctx, r)
		if err != nil {
			return false, err
		}

		return status == Running, nil
	}, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return d.docker.Network(ctx, r.Name)
}

// Connect attempts to connect the given registry to the specified docker network
func (d *DockerRegistry) Connect(ctx context.Context, r *Registry, network string) error {
	return d.docker.Connect(ctx, r.Name, network)
}

func (d *DockerRegistry) Destroy(ctx context.Context, r *Registry) error {
	cleaner := CleanDockerRegistry(r)
	return cleaner(ctx)
}

func CleanDockerRegistry(r *Registry) Cleaner {
	return func(_ context.Context) error {
		return sh.RunV("docker", "rm", "-f", r.Name)
	}
}
