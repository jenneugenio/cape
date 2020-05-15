package mage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

type Chart struct {
	Name    string
	Chart   string
	Version string // Version of the chart
	Values  string
}

type Helm struct {
	Version      *semver.Version
	Repositories map[string]string
}

func NewHelm(required string, repositories map[string]string) (*Helm, error) {
	v, err := semver.NewVersion(required)
	if err != nil {
		return nil, err
	}

	return &Helm{
		Version:      v,
		Repositories: repositories,
	}, nil
}

func (h *Helm) Name() string {
	return "helm"
}

func (h *Helm) Check(_ context.Context) error {
	out, err := sh.Output("helm", "version", "--template", "{{.Version}}")
	if err != nil {
		return err
	}

	v, err := semver.NewVersion(out)
	if err != nil {
		return fmt.Errorf("Could not parse `helm version`: %s", err.Error())
	}

	if v.LessThan(h.Version) {
		return fmt.Errorf("Please upgrade your version of Helm from %s to %s", v.String(), h.Version.String())
	}

	return nil
}

func (h *Helm) Setup(ctx context.Context) error {
	for name, repo := range h.Repositories {
		err := sh.Run("helm", "repo", "add", name, repo)
		if err != nil {
			return err
		}
	}

	return h.RepoUpdate(ctx)
}

func (h *Helm) Clean(_ context.Context) error {
	return nil
}

func (h *Helm) RepoUpdate(_ context.Context) error {
	return sh.Run("helm", "repo", "update")
}

type Release struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Revision   string `json:"revision"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
}

func (h *Helm) Deploy(ctx context.Context, c *Chart) error {
	releases, err := h.List(ctx)
	if err != nil {
		return err
	}

	var current *Release
	for _, release := range releases {
		if release.Name == c.Name {
			current = release
			break
		}
	}

	if current == nil {
		return h.Install(ctx, c)
	}

	return h.Update(ctx, c, current)
}

func (h *Helm) Install(ctx context.Context, c *Chart) error {
	args := []string{
		"install",
		"--dependency-update",
		"--version",
		c.Version,
	}

	if c.Values != "" {
		args = append(args, "--values", c.Values)
	}

	args = append(args, c.Name, c.Chart)

	return sh.RunV("helm", args...)
}

func (h *Helm) Update(ctx context.Context, c *Chart, current *Release) error {
	args := []string{
		"upgrade",
		"--atomic",
		"--cleanup-on-fail",
		"--version",
		c.Version,
	}

	if c.Values != "" {
		args = append(args, "--values", c.Values)
	}

	args = append(args, c.Name, c.Chart)
	return sh.RunV("helm", args...)
}

func (h *Helm) List(_ context.Context) ([]*Release, error) {
	out, err := sh.Output("helm", "list", "--output", "json")
	if err != nil {
		return nil, err
	}

	releases := []*Release{}
	err = json.Unmarshal([]byte(out), &releases)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal `helm list` output: %s", err.Error())
	}

	return releases, nil
}
