package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

const (
	DefaultAdminPolicy  = Label("default-admin")
	DefaultGlobalPolicy = Label("default-global")
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
	if err := p.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidPolicyCause, err)
	}

	err := p.Label.Validate()
	if err != nil {
		return errors.Wrap(InvalidPolicyCause, err)
	}

	err = p.Spec.Validate()
	if err != nil {
		return errors.Wrap(InvalidPolicyCause, err)
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
		Label:     ps.Label,
		Spec:      ps,
	}

	return policy, policy.Validate()
}

func (p *Policy) GetEncryptable() bool {
	return false
}
