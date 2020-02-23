package primitives

import (
	"time"

	"github.com/dropoutlabs/privacyai/primitives/types"
)

// TestEntity represents an internal Entity used exclusively for testing
type TestEntity struct {
	*Primitive
	Data string
}

// NewTestEntity returns a new TestEntity struct
func NewTestEntity(data string) (*TestEntity, error) {
	p, err := newPrimitive(types.Test)
	if err != nil {
		return nil, err
	}

	e := &TestEntity{
		Primitive: p,
		Data:      data,
	}

	// XXX: Static time for the purposes of testing
	e.CreatedAt = time.Unix(0, 0).UTC()
	e.UpdatedAt = time.Unix(0, 0).UTC()

	ID, err := DeriveID(e)
	if err != nil {
		return nil, err
	}

	e.ID = ID
	return e, nil
}

// TestMutableEntity represents an internal Entity used exclusively for testing
type TestMutableEntity struct {
	*Primitive
	Data string
}

// NewTestMutableEntity returns a new TestMutableEntity
func NewTestMutableEntity(data string) (*TestMutableEntity, error) {
	p, err := newPrimitive(types.TestMutable)
	if err != nil {
		return nil, err
	}

	return &TestMutableEntity{
		Primitive: p,
		Data:      data,
	}, nil
}
