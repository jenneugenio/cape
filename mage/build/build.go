package build

import (
	"context"
	"sync"

	"github.com/magefile/mage/mg"

	"github.com/capeprivacy/cape/mage"
)

func init() {
	err := mage.Tracker.Add("bin/cape")
	if err != nil {
		panic(err)
	}
}

type Build mg.Namespace

// Binary builds the Cape binary and makes it available for a use locally or to
// be packaged up for a release.
//
// Set the GOOS and GOARCH environment variables to change the target platform
// of the built artifact.
func (b Build) Binary(ctx context.Context) error {
	mg.Deps(b.Generate)

	required := []string{"git", "go"}
	err := mage.Dependencies.Run(required, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	deps, err := mage.Dependencies.Get(required)
	if err != nil {
		return err
	}

	git := deps[0].(*mage.Git)
	golang := deps[1].(*mage.Golang)

	version, err := git.Tag(ctx)
	if err != nil {
		return err
	}

	return golang.Build(ctx, version, "cmd", "bin/cape")
}

// Generate generates any required files to build the binary (GraphQL, gRPC, Protobuf)
func (b Build) Generate(ctx context.Context) error {
	err := mage.Dependencies.Run([]string{"go", "protoc"}, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	deps, err := mage.Dependencies.Get([]string{"protoc"})
	if err != nil {
		return err
	}

	gql, err := mage.NewGraphQL("coordinator/schema", "gqlgen.yml", "coordinator/graph")
	if err != nil {
		return err
	}

	errors := mage.NewSafeSlice()
	wg := &sync.WaitGroup{}

	generators := []mage.Generator{gql, deps[0].(mage.Generator)}
	for _, g := range generators {
		generator := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			errors.Add(generator.Generate(ctx))
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
