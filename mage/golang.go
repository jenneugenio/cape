package mage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	Mod     *GoMod
	Tools   *GoTools
}

func NewGolang(rootPkg string, required string) (*Golang, error) {
	v, err := semver.NewVersion(required)
	if err != nil {
		return nil, err
	}

	return &Golang{
		Version: v,
		RootPkg: rootPkg,
		Mod:     &GoMod{},
		Tools:   &GoTools{},
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
	if err := g.Mod.Setup(ctx); err != nil {
		return err
	}

	return g.Tools.Setup(ctx)
}

func (g *Golang) Clean(ctx context.Context) error {
	// Collect errors and return them as a multi error at the end!
	errors := NewErrors()

	errors.Append(g.Tools.Clean(ctx))
	errors.Append(g.Mod.Clean(ctx, g.RootPkg))
	errors.Append(sh.Run("go", "clean", "-cache", "-testcache", "-r", g.RootPkg))

	return errors.Err()
}

// Build provides functionality for building a go binary given the pkg path of
// its main and a output path for the binary.
func (g *Golang) Build(ctx context.Context, version *Version, pkg, out string) error {
	err := g.Mod.Setup(ctx)
	if err != nil {
		return err
	}

	env := map[string]string{
		"GOOS":   os.Getenv("GOOS"),
		"GOARCH": os.Getenv("GOARCH"),
	}

	pkg = filepath.Join(g.RootPkg, pkg)
	ldflags := fmt.Sprintf(`'-w -X "%s/version.Version=%s" -X "%s/version.BuildDate=%s" -s'`,
		g.RootPkg, version.Version(), g.RootPkg, version.BuildDate())
	return sh.RunWith(env, "go", "build", "-i", "-v", "-ldflags", ldflags, "-o", out, pkg)
}

// GoMod represents the `go mod` command and all of the logic associated with
// managing go modules
type GoMod struct{}

func (g *GoMod) Setup(_ context.Context) error {
	return sh.Run("go", "mod", "download")
}

func (g *GoMod) Tidy(_ context.Context) error {
	return sh.Run("go", "mod", "tidy")
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
	errors := NewErrors()
	wg := &sync.WaitGroup{}
	for _, p := range pkgs {
		pkg := p

		wg.Add(1)
		go func() {
			defer wg.Done()
			errors.Append(f(pkg))
		}()
	}

	wg.Wait()

	return errors.Err()
}
