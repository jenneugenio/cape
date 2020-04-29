package mage

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var resolverRegExp = regexp.MustCompile(`\/([a-z]+)\.graphql$`)

type GraphQL struct {
	SchemaDir    string
	ConfigPath   string
	GeneratedDir string
}

func NewGraphQL(schemaDir string, configPath string, generatedDir string) (*GraphQL, error) {
	return &GraphQL{
		SchemaDir:    schemaDir,
		ConfigPath:   configPath,
		GeneratedDir: generatedDir,
	}, nil
}

// Generate is responsible for generating all of the code used by gqlgen
func (g *GraphQL) Generate(_ context.Context) error {
	// Get a list of all the schemas
	schemas, err := filepath.Glob(filepath.Join(g.SchemaDir, "*.graphql"))
	if err != nil {
		return err
	}

	// Build a list of files that are generated via the gqlgen `go generate`
	generated := []string{
		filepath.Join(g.GeneratedDir, "generated/generated.go"),
		filepath.Join(g.GeneratedDir, "model/models_gen.go"),
	}
	for _, schema := range schemas {
		matches := resolverRegExp.FindStringSubmatch(schema)
		resolver := filepath.Join(g.GeneratedDir, fmt.Sprintf("%s.resolvers.go", matches[1]))
		generated = append(generated, resolver)
	}

	// src is a list of all the files that result in one of the generated files
	src := append(schemas, filepath.Join(g.GeneratedDir, "resolver.go"), g.ConfigPath)

	skipGeneration := true
	for _, dst := range generated {
		needsGeneration, err := target.Path(dst, src...)
		if err != nil {
			return err
		}

		if needsGeneration {
			skipGeneration = false
			break
		}
	}

	if skipGeneration {
		return nil
	}

	return sh.Run("go", "generate", filepath.Join(g.GeneratedDir, "resolver.go"))
}
