package mage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

var gitVersionRegexp = regexp.MustCompile(`git version (([0-9]+\.?)*)$`)

type Git struct {
	Version *semver.Version
}

func NewGit(required string) (*Git, error) {
	v, err := semver.NewVersion(required)
	if err != nil {
		return nil, err
	}

	return &Git{
		Version: v,
	}, nil
}

func (g *Git) Name() string {
	return "git"
}

func (g *Git) Check(_ context.Context) error {
	out, err := sh.Output("git", "--version")
	if err != nil {
		return fmt.Errorf("Could not check version of Git, please ensure it's installed: %s", err.Error())
	}

	matches := gitVersionRegexp.FindStringSubmatch(out)
	if len(matches) != 3 {
		return fmt.Errorf("Could not process output of `git --version`")
	}

	v, err := semver.NewVersion(matches[1])
	if err != nil {
		return fmt.Errorf("Could not parse `git --version`: %s", err.Error())
	}

	if v.LessThan(g.Version) {
		return fmt.Errorf("Please upgrade your version of Git to %s or greater", g.Version)
	}

	return nil
}

func (g *Git) Tag(_ context.Context) (*Version, error) {
	tags, err := sh.Output("git", "tag", "-l")
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return NewVersion("0.0.0")
	}

	HEAD, err := sh.Output("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return nil, err
	}

	return NewVersion(HEAD)
}

func (g *Git) Porcelain(_ context.Context, files []string) error {
	args := []string{"status", "--untracked-files=no", "--porcelain"}
	if len(files) > 0 {
		args = append(args, "--")
		args = append(args, files...)
	} else {
		args = append(args, ".")
	}

	out, err := sh.Output("git", args...)
	if err != nil {
		return err
	}

	if len(out) > 0 {
		return fmt.Errorf("uncommitted changes found")
	}

	return nil
}

func (g *Git) Setup(_ context.Context) error {
	return nil
}

func (g *Git) Clean(_ context.Context) error {
	return nil
}
