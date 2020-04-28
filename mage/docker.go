package mage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

var dockerVersionRegex = regexp.MustCompile(`Docker version (([0-9]+\.?)*)`)

type DockerImage struct {
	Name string
	File string
}

// Docker is a dependency check for Docker
type Docker struct {
	Version *semver.Version
}

func NewDocker(required string) (*Docker, error) {
	v, err := semver.NewVersion(required)
	if err != nil {
		return nil, err
	}

	return &Docker{
		Version: v,
	}, nil
}

// Check returns an error if Docker isn't available or the version is incorrect
func (d *Docker) Check(_ context.Context) error {
	out, err := sh.Output("docker", "-v")
	if err != nil {
		return err
	}

	matches := dockerVersionRegex.FindStringSubmatch(out)
	if len(matches) != 3 {
		return fmt.Errorf("Could not parse output of `docker -v`")
	}

	v, err := semver.NewVersion(matches[1])
	if err != nil {
		return fmt.Errorf("Could not parse `docker -v`: %s", err.Error())
	}

	if v.LessThan(d.Version) {
		return fmt.Errorf("Please upgrade your version of Docker to %s or greater", d.Version.String())
	}

	return nil
}

func (d *Docker) Build(_ context.Context, label, dockerfile string) error {
	return sh.RunV("docker", "build", "-t", label, "-f", dockerfile, ".")
}

func (d *Docker) Name() string {
	return "docker"
}

func (d *Docker) Setup(_ context.Context) error {
	return nil
}

func (d *Docker) Clean(_ context.Context) error {
	return nil
}

func CleanDockerImage(name string) Cleaner {
	return func(_ context.Context) error {
		return sh.Run("docker", "images", "rm", "-f", name)
	}
}
