package mage

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
)

var goVersionRegex = regexp.MustCompile(`go([0-9]\.[0-9]*)`)

// Golang is a dependency checker for the Go language
type Golang struct {
	Version *semver.Version
	RootPkg string
	mod     *GoMod
	tools   *GoTools
}

func NewGolang(rootPkg string, required string) (*Golang, error) {
	v, err := semver.NewVersion(required)
	if err != nil {
		return nil, err
	}

	return &Golang{
		Version: v,
		RootPkg: rootPkg,
		mod:     &GoMod{},
		tools:   &GoTools{},
	}, nil
}

// Name returns the internal "name" for the dependency
func (g *Golang) Name() string {
	return "go"
}

// Check returns an error if Go isn't available or the version is incorrect.
func (g *Golang) Check(_ context.Context) error {
	out, err := sh.Output("go", "version")
	if err != nil {
		return err
	}

	matches := goVersionRegex.FindStringSubmatch(out)
	if len(matches) != 2 {
		return fmt.Errorf("Could not process output of `go version`")
	}

	v, err := semver.NewVersion(matches[1])
	if err != nil {
		return fmt.Errorf("Could not parse `go version`: %s", err.Error())
	}

	if v.LessThan(g.Version) {
		return fmt.Errorf("Please upgrade your version of Go to %s or greater", g.Version.String())
	}

	return nil
}

func (g *Golang) Setup(ctx context.Context) error {
	if err := g.mod.Setup(ctx); err != nil {
		return err
	}

	return g.tools.Setup(ctx)
}

func (g *Golang) Clean(ctx context.Context) error {
	errors := []error{} // Collect errors and deal with them at the end
	if err := g.tools.Clean(ctx); err != nil {
		errors = append(errors, err) // Collect errors and return at the end
	}

	if err := g.mod.Clean(ctx, g.RootPkg); err != nil {
		errors = append(errors, err)
	}

	if err := sh.Run("go", "clean", "-cache", "-testcache", "-r", g.RootPkg); err != nil {
		errors = append(errors, err)
	}

	// TODO: Introduce multi-error type we can use for bucketing errors
	// together.
	if len(errors) == 0 {
		return nil
	}

	return errors[0]
}

// GoMod represents the `go mod` command and all of the logic associated with
// managing go modules
type GoMod struct{}

func (g *GoMod) Setup(_ context.Context) error {
	return sh.Run("go", "mod", "download")
}

func (g *GoMod) Clean(_ context.Context, pkg string) error {
	return sh.Run("go", "clean", "-modcache", "-r", pkg)
}

// GoTools represents the external tools that we download, install, and use as
// a part of our build pipelines.
type GoTools struct{}

func (g *GoTools) Setup(ctx context.Context) error {
	pkgs, err := g.List()
	if err != nil {
		return err
	}

	return g.run(ctx, pkgs, func(pkg string) error {
		return sh.Run("go", "install", pkg)
	})
}

func (g *GoTools) Clean(ctx context.Context) error {
	pkgs, err := g.List()
	if err != nil {
		return err
	}

	return g.run(ctx, pkgs, func(pkg string) error {
		return sh.Run("go", "clean", "-i", "-r", pkg)
	})
}

func (g *GoTools) List() ([]string, error) {
	return FindImports("tools.go")
}

func (g *GoTools) run(_ context.Context, pkgs []string, f func(pkg string) error) error {
	errors := NewSafeSlice()
	wg := &sync.WaitGroup{}
	for _, p := range pkgs {
		pkg := p

		wg.Add(1)
		go func() {
			defer wg.Done()
			errors.Add(f(pkg))
		}()
	}

	wg.Wait()
	for _, err := range errors.Get() {
		if err != nil {
			return err.(error)
		}
	}

	return nil
}
