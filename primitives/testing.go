package primitives

import (
	"github.com/dropoutlabs/privacyai/primitives/types"
)

// TestEntity represents an internal Entity used exclusively for testing
type TestEntity struct {
	Primitive
	Data string
}

// GetType returns the Type for this Entity
func (t *TestEntity) GetType() types.Type {
	return types.Test
}

// TestingMutableEntity represents an internal Entity used exclusively for testing
type TestingMutableEntity struct {
	Primitive
	Data string
}

// GetType returns the Type for this Entity
func (t *TestingMutableEntity) GetType() types.Type {
	return types.TestMutable
}
