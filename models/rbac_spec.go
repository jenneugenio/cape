package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"sigs.k8s.io/yaml"
)

// RBACRule represents the rules that make up a policy
type RBACRule struct {
	Target Target     `json:"target"`
	Action RBACAction `json:"action"`
	Effect Effect     `json:"effect"`
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (r *RBACRule) UnmarshalGQL(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(t, r); err != nil {
			return err
		}

		return r.Validate()
	default:
		return errors.New("unable to unmarshal gql policy spec")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (r RBACRule) MarshalGQL(w io.Writer) {
	json, err := json.Marshal(r)
	if err != nil {
		fmt.Fprint(w, strconv.Quote(err.Error()))
		return
	}

	fmt.Fprint(w, string(json))
}

// Validate validates that the rule arguments are valid
func (r *RBACRule) Validate() error {
	err := r.Target.Validate()
	if err != nil {
		return fmt.Errorf("rule has invalid target: %w", err)
	}

	return nil
}

// RBACSpec defines the policy (e.g. the yaml file)
type RBACSpec struct {
	Label Label       `json:"label"`
	Rules []*RBACRule `json:"rules"`
}

// Validate that the policy spec is valid
func (ps *RBACSpec) Validate() error {
	for _, r := range ps.Rules {
		err := r.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (ps *RBACSpec) UnmarshalGQL(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(t, ps); err != nil {
			return err
		}

		return ps.Validate()
	default:
		return errors.New("Unable to unmarshal gql policy spec")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (ps RBACSpec) MarshalGQL(w io.Writer) {
	json, err := json.Marshal(ps)
	if err != nil {
		fmt.Fprint(w, strconv.Quote(err.Error()))
		return
	}

	fmt.Fprint(w, string(json))
}

// ToBytes writes the policy spec to bytes representing the policy file
func (ps *RBACSpec) ToBytes() ([]byte, error) {
	return yaml.Marshal(ps)
}

// ParseRBACSpec reads a policy spec file and returns a policy spec object
func ParseRBACSpec(data []byte) (*RBACSpec, error) {
	var ps RBACSpec
	err := yaml.Unmarshal(data, &ps)
	if err != nil {
		return nil, err
	}

	return &ps, ps.Validate()
}
