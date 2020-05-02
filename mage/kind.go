package mage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

var kubectlVersionRegex = regexp.MustCompile(`Client Version: v(([0-9]+\.?)*)`)
var kindVersionRegex = regexp.MustCompile(`kind v(([0-9+\.?])*)`)

type Cluster struct {
	Name       string
	ConfigFile string
}

// Kind represents a dependency on KinD (Kubernetes in Docker) which is used to
// create and manage our local clusters.
//
// Kind is installed via go.mod and specified as a dependency inside of tools.go
type Kind struct {
	Version     *semver.Version
	docker      *Docker
	k8sVersion  *semver.Version
	cfgTemplate *template.Template
}

func NewKind(docker *Docker, version string, k8sVersion string) (*Kind, error) {
	if docker == nil {
		return nil, errors.New("A docker registry must be provided")
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	k8sSemver, err := semver.NewVersion(k8sVersion)
	if err != nil {
		return nil, err
	}

	t, err := template.ParseFiles("mage/config/kind.yaml.template")
	if err != nil {
		return nil, err
	}

	return &Kind{
		Version:     v,
		docker:      docker,
		k8sVersion:  k8sSemver,
		cfgTemplate: t,
	}, nil
}

func (k *Kind) Check(ctx context.Context) error {
	kindOut, err := sh.Output("kind", "version")
	if err != nil {
		return err
	}

	kvMatches := kindVersionRegex.FindStringSubmatch(kindOut)
	if len(kvMatches) != 3 {
		return fmt.Errorf("Could not parse output of `kind version`")
	}

	kv, err := semver.NewVersion(kvMatches[1])
	if err != nil {
		return fmt.Errorf("Could not parse output of `kind version`: %s", err.Error())
	}

	if kv.LessThan(k.Version) {
		return fmt.Errorf("Your version of kind is out of date, please upgrade to %s", k.Version.String())
	}

	kubectlOut, err := sh.Output("kubectl", "version", "--short", "--client")
	if err != nil {
		return err
	}

	matches := kubectlVersionRegex.FindStringSubmatch(kubectlOut)
	if len(matches) != 3 {
		return fmt.Errorf("Could not parse output of `kubectl version --short --client`")
	}

	v, err := semver.NewVersion(matches[1])
	if err != nil {
		return fmt.Errorf("Could not parse output of `kubectl version --short --client`: %s", err.Error())
	}

	if v.LessThan(k.k8sVersion) {
		return fmt.Errorf("Please upgrade your version of kubectl from %s to %s or greater", v.String(), k.k8sVersion.String())
	}

	return k.docker.Check(ctx)
}

func (k *Kind) Name() string {
	return "kind"
}

func (k *Kind) Setup(_ context.Context) error {
	return nil
}

func (k *Kind) Clean(_ context.Context) error {
	return nil
}

// Create creates a configured KinD cluster with a local docker registry for
// making docker files available to the kubernetes cluster.
func (k *Kind) Create(ctx context.Context, cluster *Cluster, registry *Registry) error {
	state, err := k.Status(ctx, cluster)
	if err != nil {
		return err
	}

	if state == Running {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFile := filepath.Join(wd, filepath.Clean(cluster.ConfigFile))
	dir := filepath.Dir(cluster.ConfigFile)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("Could not create build directory %s: %s", dir, err.Error())
	}

	f, err := os.Create(cfgFile)
	if err != nil {
		return fmt.Errorf("Could not create config file: %s: %s", cfgFile, err.Error())
	}
	defer f.Close()

	err = k.cfgTemplate.Execute(f, map[string]interface{}{
		"Registry": registry,
	})
	if err != nil {
		return fmt.Errorf("Error writing config to %s: %s", cfgFile, err.Error())
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("Could not sync config file to disk %s: %s", cfgFile, err.Error())
	}

	nameArg := fmt.Sprintf("--name=%s", cluster.Name)
	cfgArg := fmt.Sprintf("--config=%s", cfgFile)
	return sh.Run("kind", "create", "cluster", nameArg, cfgArg)
}

func (k *Kind) Status(ctx context.Context, cluster *Cluster) (Status, error) {
	// We don't use sh.Output here -we want to stop forwarding stderr out.
	out := &strings.Builder{}
	_, err := sh.Exec(nil, out, nil, "kind", "get", "clusters")
	if err != nil {
		return Dead, err
	}

	clusters := strings.Split(out.String(), "\n")
	for _, c := range clusters {
		if c == cluster.Name {
			return Running, nil
		}
	}

	return Dead, nil
}

func (k *Kind) Destroy(ctx context.Context, cluster *Cluster) error {
	cleaner := CleanKind(cluster)
	return cleaner(ctx)
}

func CleanKind(cluster *Cluster) Cleaner {
	return func(_ context.Context) error {
		errors := NewErrors()
		errors.Append(os.Remove(cluster.ConfigFile))

		nameArg := fmt.Sprintf("--name=%s", cluster.Name)
		errors.Append(sh.Run("kind", "delete", "cluster", nameArg))

		return errors.Err()
	}
}
