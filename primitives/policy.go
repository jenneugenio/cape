package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Policy is a single defined policy
type Policy struct {
	*database.Primitive
	Label Label
}

// GetType returns the type for this entity
func (p *Policy) GetType() types.Type {
	return PolicyType
}

// NewPolicy returns a mutable policy struct
func NewPolicy(label Label) (*Policy, error) {
	p, err := database.NewPrimitive(PolicyType)
	if err != nil {
		return nil, err
	}

	return &Policy{
		Primitive: p,
		Label:     label,
	}, nil
}
