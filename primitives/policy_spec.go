package primitives

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"io"
	"sigs.k8s.io/yaml"
	"strconv"
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

// MarshalPolicySpec gql implementation
func MarshalPolicySpec(ps PolicySpec) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		yaml, err := yaml.Marshal(ps)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(w, strconv.Quote(string(yaml)))
	})
}

// UnmarshalPolicySpec gql implementation
func UnmarshalPolicySpec(v interface{}) (PolicySpec, error) {
	var ps PolicySpec
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return ps, err
	}

	err = yaml.Unmarshal(bytes, &ps)
	if err != nil {
		return ps, err
	}

	return PolicySpec{}, nil
}

// ToBytes writes the policy spec to bytes representing the policy file
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
