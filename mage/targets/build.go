package targets

import (
	"context"
	"fmt"
	"sync"

	"github.com/magefile/mage/mg"

	"github.com/capeprivacy/cape/mage"
)

var dockerImages = []mage.DockerImage{
	{
		Name: "capeprivacy/base",
		File: "dockerfiles/Dockerfile",
	},
	{
		Name: "capeprivacy/cape",
		File: "dockerfiles/Dockerfile.cape",
	},
	{
		Name: "capeprivacy/update",
		File: "dockerfiles/Dockerfile.update",
	},
	{
		Name: "capeprivacy/customer_seed",
		File: "tools/seed/Dockerfile.customer",
	},
}

func init() {
	err := mage.Tracker.Add("file://bin/cape", mage.CleanFile("bin/cape"))
	if err != nil {
		panic(err)
	}

	for _, image := range dockerImages {
		err := mage.Tracker.Add("docker://"+image.Name, mage.CleanDockerImage(image.Name))
		if err != nil {
			panic(err)
		}
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

// Docker builds all of the cape docker containers
func (b Build) Docker(ctx context.Context) error {
	required := []string{"docker"}
	err := mage.Dependencies.Run([]string{"docker"}, func(d mage.Dependency) error {
		return d.Check(ctx)
	})
	if err != nil {
		return err
	}

	deps, err := mage.Dependencies.Get(required)
	if err != nil {
		return err
	}

	docker := deps[0].(*mage.Docker)
	for _, image := range dockerImages {
		err := docker.Build(ctx, fmt.Sprintf("%s:latest", image.Name), image.File)
		if err != nil {
			return err
		}
	}

	return nil
}
