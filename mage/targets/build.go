package targets

import (
	"context"
	"sync"

	"github.com/magefile/mage/mg"

	"github.com/capeprivacy/cape/mage"
)

var dockerImages = []*mage.DockerImage{
	{
		Name: "capeprivacy/base",
		Tag:  "latest",
		File: "dockerfiles/Dockerfile",
	},
	{
		Name: "capeprivacy/cape",
		Tag:  "latest",
		File: "dockerfiles/Dockerfile.cape",
		Args: func(ctx context.Context) (map[string]string, error) {
			version, err := getVersion(ctx)
			if err != nil {
				return nil, err
			}

			return map[string]string{
				"VERSION": version.Version(),
			}, nil
		},
	},
	{
		Name: "capeprivacy/update",
		Tag:  "latest",
		File: "dockerfiles/Dockerfile.update",
	},
	{
		Name: "capeprivacy/customer_seed",
		Tag:  "latest",
		File: "tools/seed/Dockerfile.customer",
	},
}

func init() {
	err := mage.Tracker.Add("file://bin/cape", mage.CleanFile("bin/cape"))
	if err != nil {
		panic(err)
	}

	for _, image := range dockerImages {
		err := mage.Tracker.Add("docker://"+image.Name, mage.CleanDockerImage(image))
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

	required := []string{"go"}
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

	golang := deps[0].(*mage.Golang)

	// We optionally support specifying the version via an env variable. This
	// gets around us having to pull the version from Git which is useful when
	// we're building Cape inside of a container without a checkout of git.
	version, err := getVersion(ctx)
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

	pkger, err := mage.NewPkger()
	if err != nil {
		return err
	}

	errors := mage.NewErrors()
	wg := &sync.WaitGroup{}

	generators := []mage.Generator{gql, pkger, deps[0].(mage.Generator)}
	for _, g := range generators {
		generator := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			errors.Append(generator.Generate(ctx))
		}()
	}

	wg.Wait()

	return errors.Err()
}

// Docker builds all of the cape docker containers
func (b Build) Docker(ctx context.Context) error {
	required := []string{"docker"}
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

	docker := deps[0].(*mage.Docker)
	for _, image := range dockerImages {
		err := docker.Build(ctx, image)
		if err != nil {
			return err
		}
	}

	return nil
}
