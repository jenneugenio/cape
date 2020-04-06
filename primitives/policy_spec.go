package primitives

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"sigs.k8s.io/yaml"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/mitchellh/mapstructure"
)

// PolicyVersion represents which version of policy we are using
type PolicyVersion uint8

// Where expresses a condition that says if a rule should apply to certain data
type Where map[string]string

// Operation represents an operation in a query (e.g. equality check)
type Operation string

const (
	Eq  Operation = "="
	Neq Operation = "!="
)

type RuleType int

const (
	FieldRule RuleType = iota
	WhereRule
)

func (r *Rule) Type() RuleType {
	if len(r.Fields) > 0 {
		return FieldRule
	}

	return WhereRule
}

// Rule represents the rules that make up a policy
type Rule struct {
	Target Target  `json:"target"`
	Action Action  `json:"action"`
	Effect Effect  `json:"effect"`
	Fields []Field `json:"fields,omitempty"`
	Where  []Where `json:"where,omitempty"`
	Sudo   bool    `json:"sudo"`
}

// PolicySpec defines the policy (e.g. the yaml file)
type PolicySpec struct {
	Version PolicyVersion `json:"version"`
	Label   Label         `json:"label"`
	Rules   []*Rule       `json:"rules"`
}

// Validate that the policy spec is valid
func (ps *PolicySpec) Validate() error {
	for _, r := range ps.Rules {
		err := r.Target.Validate()
		if err != nil {
			return errors.New(InvalidPolicySpecCause, "Invalid rule: %e", err)
		}

		if len(r.Fields) > 0 && len(r.Where) > 0 {
			return errors.New(InvalidPolicySpecCause, "Fields & Where cannot be specified on the same rule")
		}
	}

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (ps *PolicySpec) UnmarshalGQL(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		return mapstructure.Decode(t, ps)
	default:
		return errors.New(InvalidPolicySpecCause, "Unable to unmarshal gql policy spec")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (ps PolicySpec) MarshalGQL(w io.Writer) {
	json, err := json.Marshal(ps)
	if err != nil {
		fmt.Fprint(w, strconv.Quote(err.Error()))
	}

	fmt.Fprint(w, string(json))
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
