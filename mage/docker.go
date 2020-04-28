package mage

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

var dockerVersionRegex = regexp.MustCompile(`Docker version (([0-9]+\.?)*)`)

type DockerImage struct {
	Name string
	File string
	Tag  string
}

func (i *DockerImage) String() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
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
	if len(os.Getenv("SKIP_DOCKER_CHECK")) > 0 {
		return nil
	}

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
		return fmt.Errorf("Please upgrade your version of Docker from %s to %s or greater", v.String(), d.Version.String())
	}

	return nil
}

func (d *Docker) Build(_ context.Context, label, dockerfile string, args map[string]string) error {
	cmd := []string{"build", "-t", label, "-f", dockerfile}
	if len(args) > 0 {
		for key, value := range args {
			cmd = append(cmd, "--build-arg", fmt.Sprintf("\"%s=%s\"", key, value))
		}
	}

	cmd = append(cmd, ".")
	return sh.RunV("docker", cmd...)
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

func CleanDockerImage(image *DockerImage) Cleaner {
	return func(_ context.Context) error {
		return sh.Run("docker", "rm", "-f", image.String())
	}
}
