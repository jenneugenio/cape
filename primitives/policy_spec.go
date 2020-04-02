package primitives

import (
	"fmt"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"sigs.k8s.io/yaml"
)

// Version represents which version of policy we are using
type PolicyVersion uint8

// Where expresses a condition that says if a rule should apply to certain data
type Where map[string]string

// Rule represents the rules that make up a policy
type Rule struct {
	Target Target  `json:"target"`
	Action Action  `json:"action"`
	Effect Effect  `json:"effect"`
	Fields []Field `json:"fields"`
	Where  []Where `json:"where,omitempty"`
}

// PolicySpec defines the policy (e.g. the yaml file)
type PolicySpec struct {
	Version PolicyVersion `json:"version"`
	Label   Label         `json:"label"`
	Rules   []Rule        `json:"rules"`
}

// Validate that the policy spec is valid
func (ps *PolicySpec) Validate() error {
	for _, r := range ps.Rules {
		err := r.Target.Validate()
		if err != nil {
			return errors.New(InvalidPolicySpecCause, fmt.Sprintf("Invalid rule: %e", err))
		}
	}

	return nil
}

func (ps *PolicySpec) ToBytes() ([]byte, error) {
	return yaml.Marshal(ps)
}

// ParsePolicySpec reads a policy spec file and returns a policy spec object
func ParsePolicySpec(data []byte) (*PolicySpec, error) {
	var ps PolicySpec
	err := yaml.Unmarshal(data, &ps)
	if err != nil {
		return nil, err
	}

	return &ps, ps.Validate()
}
