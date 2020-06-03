package primitives

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"sigs.k8s.io/yaml"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/transformations"
)

// PolicyVersion represents which version of policy we are using
type PolicyVersion uint8

// Where expresses a condition that says if a rule should apply to certain data
type Where map[string]string

// Operation represents an operation in a query (e.g. equality check)
type Operation string

const (
	// Eq is the equality operator
	Eq Operation = "="

	// Neq is the not equal operator
	Neq Operation = "!="
)

// RuleType is the type of a rule
type RuleType int

const (
	// FieldRule is a rule that contains a field clause it cannot
	// be specified if a where clause is specified
	FieldRule RuleType = iota

	// WhereRule is a rule that contains a where clause it cannot
	// be specified if a field clause is specified
	WhereRule
)

// Rule represents the rules that make up a policy
type Rule struct {
	Target          Target            `json:"target"`
	Action          Action            `json:"action"`
	Effect          Effect            `json:"effect"`
	Fields          []Field           `json:"fields,omitempty"`
	Where           []Where           `json:"where,omitempty"`
	Transformations []*Transformation `json:"transformations,omitempty"`
	Sudo            bool              `json:"sudo"`
}

// Validate validates that the rule arguments are valid
func (r *Rule) Validate() error {
	err := r.Target.Validate()
	if err != nil {
		return errors.Wrap(InvalidPolicySpecCause, err)
	}

	if r.Target.Type() == Records {
		if len(r.Fields) > 0 && len(r.Where) > 0 {
			return errors.New(InvalidPolicySpecCause, "Fields & Where cannot be specified on the same rule")
		}
	} else if len(r.Fields) > 0 || len(r.Where) > 0 {
		return errors.New(InvalidPolicySpecCause, "Fields & Where cannot be specified for a "+
			"non records policy type %s", r.Target.Type())
	}

	if r.Effect == Deny && len(r.Transformations) > 0 {
		return errors.New(InvalidPolicySpecCause, "Deny rules cannot have transformations")
	}

	for _, transform := range r.Transformations {
		err := transform.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Type returns the type of a rule
func (r *Rule) Type() RuleType {
	if len(r.Fields) > 0 {
		return FieldRule
	}

	return WhereRule
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
		err := r.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (ps *PolicySpec) UnmarshalGQL(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(t, ps); err != nil {
			return err
		}

		return ps.Validate()
	default:
		return errors.New(InvalidPolicySpecCause, "Unable to unmarshal gql policy spec")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (ps PolicySpec) MarshalGQL(w io.Writer) {
	json, err := json.Marshal(ps)
	if err != nil {
		fmt.Fprint(w, strconv.Quote(err.Error()))
		return
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

// Transformation represents a transform in the policy spec
type Transformation struct {
	Field    Field                     `json:"field"`
	Function string                    `json:"function"`
	Args     transformations.Args      `json:"args"`
	Where    transformations.Condition `json:"where,omitempty"`
}

// Validate that the policy spec is valids
func (t *Transformation) Validate() error {
	err := t.Field.Validate()
	if err != nil {
		return err
	}

	err = t.Where.Validate()
	if err != nil {
		return err
	}

	ctor := transformations.Get(t.Function)

	transform, err := ctor(t.Field.String())
	if err != nil {
		return err
	}

	return transform.Validate(t.Args)
}
