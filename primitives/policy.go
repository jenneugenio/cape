package primitives

import (
	"io/ioutil"

	"sigs.k8s.io/yaml"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Policy is a single defined policy
type Policy struct {
	*database.Primitive
	Label Label `json:"label"`
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

// Validate validates the policy
func (p *Policy) Validate() error {
	return nil
}

// ParsePolicy parses a polic file into policy p
func ParsePolicy(filePath string) (*Policy, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	policy := &Policy{}
	err = yaml.Unmarshal(b, policy)
	if err != nil {
		return nil, err
	}

	return policy, policy.Validate()
}
