package primitives

import (
	"github.com/dropoutlabs/privacyai/primitives/types"
)

// Policy is a single defined policy
// TODO -- write this
type Policy struct {
	*Primitive
}

// GetType returns the type for this entity
func (p *Policy) GetType() types.Type {
	return PolicyType
}

// NewPolicy returns a mutable policy struct
func NewPolicy() (*Policy, error) {
	p, err := newPrimitive(PolicyType)
	if err != nil {
		return nil, err
	}

	return &Policy{
		Primitive: p,
	}, nil
}
