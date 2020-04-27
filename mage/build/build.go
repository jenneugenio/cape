// +build mage
package build

import (
	"context"
	"sync"

	"github.com/magefile/mage/mg"

	"github.com/capeprivacy/cape/mage"
)

type Build mg.Namespace

// Binary builds the Cape binary and makes it available for a use locally or to
// be packaged up for a release.
func (b Build) Binary(_ context.Context) error {
	mg.Deps(b.Generate)
	return nil
}

// Generate generates any required files to build the binary (e.g. GraphQL,
// gRPC, Protobuf)
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
