package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// Policy is a single defined policy
type Policy struct {
	*database.Primitive
	Label Label       `json:"label"`
	Spec  *PolicySpec `json:"spec"`
}

// GetType returns the type for this entity
func (p *Policy) GetType() types.Type {
	return PolicyType
}

// Validate that the policy is valid
func (p *Policy) Validate() error {
	err := p.Label.Validate()
	if err != nil {
		return errors.New(InvalidPolicyCause, "invalid policy: %e", err)
	}

	err = p.Spec.Validate()
	if err != nil {
		return errors.New(InvalidPolicyCause, "invalid policy: %e", err)
	}

	return nil
}

// NewPolicy returns a mutable policy struct
func NewPolicy(label Label, spec *PolicySpec) (*Policy, error) {
	p, err := database.NewPrimitive(PolicyType)
	if err != nil {
		return nil, err
	}

	policy := &Policy{
		Primitive: p,
		Label:     label,
		Spec:      spec,
	}

	return policy, policy.Validate()
}

// ParsePolicy can convert a yaml document into a Policy
func ParsePolicy(data []byte) (*Policy, error) {
	ps, err := ParsePolicySpec(data)
	if err != nil {
		return nil, err
	}

	p, err := database.NewPrimitive(PolicyType)
	if err != nil {
		return nil, err
	}

	policy := &Policy{
		Primitive: p,
		Spec:      ps,
	}

	return policy, policy.Validate()
}
