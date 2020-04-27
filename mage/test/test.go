package test

import (
	"context"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/capeprivacy/cape/mage"
	"github.com/capeprivacy/cape/mage/build"
)

type Test mg.Namespace

// All runs the entire test suite (lint, unit and integration tests)
func (t Test) All(ctx context.Context) error {
	mg.Deps(t.Lint, t.Integration)
	return nil
}

// Lint runs only the code linting portion of the test suite
func (t Test) Lint(ctx context.Context) error {
	mg.Deps(build.Build.Generate)

	required := []string{"go"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	return sh.RunV("golangci-lint", "run")
}

// Unit runs the unit test portion of the test suite (does not require postgres)
func (t Test) Unit(ctx context.Context) error {
	mg.Deps(build.Build.Generate)

	required := []string{"go"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	return sh.RunV("go", "test", "-v", "./...")
}

// Integration runs the integration portion of the test suite (requires postgres)
func (t Test) Integration(_ context.Context) error {
	mg.Deps(build.Build.Generate)

	env := mage.Env{
		"CAPE_DB_URL":             "postgres://postgres:dev@localhost:5432/postgres?sslmode=disable",
		"CAPE_DB_MIGRATIONS":      "migrations",
		"CAPE_DB_TEST_MIGRATIONS": "database/dbtest/migrations",
		"CAPE_DB_SEED_MIGRATIONS": "tools/seed",
	}
	env.Source()

	_, err := sh.Exec(env, os.Stdout, os.Stderr, "go", "test", "-v", "./...", "-tags=integration")
	return err
}
