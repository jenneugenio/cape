package targets

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/capeprivacy/cape/mage"
)

type Test mg.Namespace

// All runs the entire test suite (lint, unit and integration tests)
func (t Test) All(ctx context.Context) error {
	// TODO: Add helm lint -- so we lint our helm packs
	mg.SerialCtxDeps(ctx, t.Lint, t.Integration, Build.Docker)
	return nil
}

// Lint runs only the code linting portion of the test suite
func (t Test) Lint(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, Build.Generate)

	required := []string{"go"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	return sh.RunV("golangci-lint", "run")
}

// Unit runs the unit test portion of the test suite (does not require Postgres)
func (t Test) Unit(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, Build.Generate)

	required := []string{"go"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	return sh.RunV("go", "test", "-v", "./...")
}

// Integration runs the integration portion of the test suite (requires Postgres)
func (t Test) Integration(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, Build.Generate)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	args := []string{"test"}

	if _, quiet := os.LookupEnv("CAPE_TEST_QUIET"); !quiet {
		args = append(args, "-v")
	}

	file := os.Getenv("CAPE_TEST_FILE")
	if file == "" {
		file = "./..."
	}
	args = append(args, file)

	args = append(args, "-tags=integration")

	coverage := os.Getenv("CAPE_DO_COVERAGE")
	if coverage != "" {
		args = append(args, "-coverprofile=coverage.txt", "-covermode=atomic")
	}

	env := mage.Env{
		"CAPE_DB_URL":             "postgres://postgres:dev@localhost:5432/postgres?sslmode=disable",
		"CAPE_DB_MIGRATIONS":      filepath.Join(wd, "coordinator/migrations"),
		"CAPE_DB_TEST_MIGRATIONS": filepath.Join(wd, "coordinator/database/dbtest/migrations"),
		"CAPE_DB_SEED_MIGRATIONS": filepath.Join(wd, "tools/seed"),
	}
	env.Source()

	_, err = sh.Exec(env, os.Stdout, os.Stderr, "go", args...)
	return err
}

// CI runs the full test suite that is run during continuous integration
func (t Test) CI(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, t.Lint, t.Integration, Build.Docker, t.Tidy)
	return nil
}

// Tidy runs `go mod tidy` and then checks for changes in `go.mod` and `go.sum`
// to ensure that go.mod and go.sum are up to date.
func (t Test) Tidy(ctx context.Context) error {
	required := []string{"git", "go"}
	deps, err := mage.Dependencies.Get(required)
	if err != nil {
		return err
	}

	git := deps[0].(*mage.Git)
	golang := deps[1].(*mage.Golang)

	err = git.Check(ctx)
	if err != nil {
		return err
	}

	err = golang.Check(ctx)
	if err != nil {
		return err
	}

	err = golang.Mod.Tidy(ctx)
	if err != nil {
		return err
	}

	files := []string{"go.mod", "go.sum"}
	err = git.Porcelain(ctx, files)
	if err != nil {
		return fmt.Errorf("Files %s have been modified: %s", files, err)
	}

	return nil
}
