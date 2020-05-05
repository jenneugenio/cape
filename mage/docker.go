package mage

import (
	"context"
	"encoding/json"
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
	Args func(context.Context) (map[string]string, error)
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

	out, err = sh.Output("docker", "info", "-f", "{{json .ServerErrors}}")
	if err != nil {
		return fmt.Errorf("Could not run `docker info`: %s", err.Error())
	}

	if out != "null" {
		return fmt.Errorf("Encountered docker server errors: %s", out)
	}

	return nil
}

func (d *Docker) Build(ctx context.Context, image *DockerImage) error {
	if len(os.Getenv("CAPE_SKIP_DOCKER_BUILD")) > 0 {
		return nil
	}

	cmd := []string{"build", "-t", image.String(), "-f", image.File}
	if image.Args != nil {
		args, err := image.Args(ctx)
		if err != nil {
			return err
		}

		for key, value := range args {
			cmd = append(cmd, "--build-arg", fmt.Sprintf("%s=%s", key, value))
		}
	}

	cmd = append(cmd, ".")
	return sh.RunV("docker", cmd...)
}

type NetworkSettings struct {
	Networks map[string]struct {
		IPAddress string `json:"IPAddress"`
	} `json:"Networks"`
	IPAddress string `json:"IPAddress"`
}

func (d *Docker) Network(_ context.Context, label string) (*NetworkSettings, error) {
	out, err := sh.Output("docker", "inspect", "-f", "{{json .NetworkSettings}}", label)
	if err != nil {
		return nil, err
	}

	settings := &NetworkSettings{}
	err = json.Unmarshal([]byte(out), settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (d *Docker) Status(ctx context.Context, name string) (Status, error) {
	nameFilter := fmt.Sprintf("name=%s", name)
	out, err := sh.Output("docker", "ps", "--filter", nameFilter, "--filter", "status=running", "--format={{.}}")
	if err != nil {
		return Unknown, err
	}

	if len(out) == 0 {
		return Unknown, nil
	}

	return Running, nil
}

func (d *Docker) Connect(ctx context.Context, name, network string) error {
	return sh.Run("docker", "network", "connect", network, name)
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
