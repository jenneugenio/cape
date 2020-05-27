package mage

import (
	"context"

	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

type Pkger struct {
}

func NewPkger() (*Pkger, error) {
	return &Pkger{}, nil
}

// Generate is responsible for generating all of the code used by pkger
func (p *Pkger) Generate(_ context.Context) error {
	src := []string{
		"primitives/policies/default/default-admin.yaml",
		"primitives/policies/default/default-global.yaml",
		"primitives/policies/default/default-data-connector.yaml",
		"primitives/policies/default/default-worker.yaml",
	}

	dst := "pkged.go"

	needsGeneration, err := target.Path(dst, src...)
	if err != nil {
		return err
	}

	if !needsGeneration {
		return nil
	}

	return sh.Run("pkger")
}
