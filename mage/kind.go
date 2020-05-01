package mage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

var kubectlVersionRegex = regexp.MustCompile(`Client Version: v(([0-9]+\.?)*)`)

type Cluster struct {
	Name       string
	ConfigFile string
}

type Kind struct {
	docker     *Docker
	k8sVersion *semver.Version
}

func NewKind(docker *Docker, k8sVersion string) (*Kind, error) {
	if docker == nil {
		return nil, errors.New("Docker must be provided")
	}

	v, err := semver.NewVersion(k8sVersion)
	if err != nil {
		return nil, err
	}

	return &Kind{
		docker:     docker,
		k8sVersion: v,
	}, nil
}

func (k *Kind) Check(ctx context.Context) error {
	out, err := sh.Output("kubectl", "version", "--short", "--client")
	if err != nil {
		return err
	}

	matches := kubectlVersionRegex.FindStringSubmatch(out)
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

func (k *Kind) Create(_ context.Context, cluster *Cluster) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	nameArg := fmt.Sprintf("--name=%s", cluster.Name)
	cfgFile := fmt.Sprintf("--config=%s", filepath.Join(wd, cluster.ConfigFile))

	return sh.RunV("kind", "create", "cluster", nameArg, cfgFile)
}

func (k *Kind) Status(ctx context.Context, cluster *Cluster) error {
	// TODO Add check or graceful support for 'cluster-info' failing in the
	// case where the cluster does not yet exist.
	return sh.RunV("kubectl", "cluster-info", "--context", fmt.Sprintf("kind-%s", cluster.Name))
}

func (k *Kind) Destroy(ctx context.Context, cluster *Cluster) error {
	cleaner := CleanKind(cluster)
	return cleaner(ctx)
}

func CleanKind(cluster *Cluster) Cleaner {
	return func(_ context.Context) error {
		nameArg := fmt.Sprintf("--name=%s", cluster.Name)
		return sh.RunV("kind", "delete", "cluster", nameArg)
	}
}
