package mage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

var kubectlVersionRegex = regexp.MustCompile(`Client Version: v(([0-9]+\.?)*)`)

type Kubectl struct {
	Version *semver.Version
}

func NewKubectl(required string) (*Kubectl, error) {
	v, err := semver.NewVersion(required)
	if err != nil {
		return nil, err
	}

	return &Kubectl{
		Version: v,
	}, nil
}

func MustKubectl(required string) *Kubectl {
	return &Kubectl{Version: semver.MustParse(required)}
}

func (k *Kubectl) Check(_ context.Context) error {
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

	if v.LessThan(k.Version) {
		return fmt.Errorf("Please upgrade your version of kubectl from %s to %s or greater", v.String(), k.Version.String())
	}

	return nil
}
