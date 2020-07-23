package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/manifoldco/go-base64"
	"github.com/mitchellh/mapstructure"
	"sigs.k8s.io/yaml"
)

func NewPolicy(label Label, rules []Rule) Policy {
	return Policy{
		ID:        NewID(),
		Version:   modelVersion,
		Label:     label,
		Rules:     rules,
		CreatedAt: now(),
	}
}

type Policy struct {
	ID              string                `json:"id"`
	Version         uint8                 `json:"version"`
	Label           Label                 `json:"label"`
	Transformations []NamedTransformation `json:"transformations"`

	Rules []Rule `json:"rules"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ParsePolicy(data []byte) (*Policy, error) {
	var p Policy

	err := yaml.Unmarshal(data, &p, func(dec *json.Decoder) *json.Decoder {
		dec.DisallowUnknownFields()
		return dec
	})
	if err != nil {
		return nil, err
	}

	return &p, nil
}

type Rule struct {
	Match   Match    `json:"match"`
	Actions []Action `json:"actions"`
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (r *Rule) UnmarshalGQL(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(t, r); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unable to unmarshal rule")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (r Rule) MarshalGQL(w io.Writer) {
	json, err := json.Marshal(r)
	if err != nil {
		fmt.Fprint(w, strconv.Quote(err.Error()))
		return
	}

	fmt.Fprint(w, string(json))
}

type Match struct {
	Name string `json:"name"`
}

type Action struct {
	Transform Transformation `json:"transform"`
}

type Transformation map[string]interface{}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *Transformation) UnmarshalGQL(v interface{}) error {
	switch val := v.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(val, t); err != nil {
			return err
		}

		return nil
	default:
		return errors.New("unable to unmarshal gql policy spec")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (t Transformation) MarshalGQL(w io.Writer) {
	json, err := json.Marshal(t)
	if err != nil {
		fmt.Fprint(w, strconv.Quote(err.Error()))
		return
	}

	fmt.Fprint(w, string(json))
}

type NamedTransformation struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Args map[string]interface{}
}

func (n NamedTransformation) MarshalJSON() ([]byte, error) {
	val := make(map[string]interface{})
	val["name"] = n.Name
	val["type"] = n.Type

	for key, arg := range n.Args {
		val[key] = arg
	}

	return json.Marshal(val)
}

func (n *NamedTransformation) UnmarshalJSON(data []byte) error {
	var transformMap map[string]interface{}
	err := json.Unmarshal(data, &transformMap)
	if err != nil {
		return err
	}

	n.Name = transformMap["name"].(string)
	n.Type = transformMap["type"].(string)

	delete(transformMap, "name")
	delete(transformMap, "type")

	n.Args = transformMap

	return findSecretArgs(n.Args)
}

// UnmarshalGQL implements the graphql.Unmarshaler interfacemin
func (n *NamedTransformation) UnmarshalGQL(v interface{}) error {
	switch val := v.(type) {
	case map[string]interface{}:
		n.Name = val["name"].(string)
		n.Type = val["type"].(string)

		delete(val, "name")
		delete(val, "type")

		n.Args = val
		return findSecretArgs(n.Args)
	default:
		return errors.New("unable to unmarshal gql policy spec")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (n NamedTransformation) MarshalGQL(w io.Writer) {
	json, err := json.Marshal(n)
	if err != nil {
		fmt.Fprint(w, strconv.Quote(err.Error()))
		return
	}

	fmt.Fprint(w, string(json))
}

type SecretArg struct {
	Type  string        `json:"type"`
	Name  string        `json:"name"`
	Value *base64.Value `json:"value"`
}

// findSecretArgs is used by UnmarshalGQL and UnmarshalJSON to find secret args
// amongst the other generic args. These can then later be encrypted or treated
// differently.
func findSecretArgs(args map[string]interface{}) error {
	// Need to convert generic args to SecretArgs if they exist
	for key, arg := range args {
		argMap, ok := arg.(map[string]interface{})
		if !ok {
			continue
		}

		sec := SecretArg{
			Type: argMap["type"].(string),
			Name: argMap["name"].(string),
		}

		val, ok := argMap["value"].(string)
		if ok {
			bVal, err := base64.NewFromString(val)
			if err != nil {
				return err
			}
			sec.Value = bVal
		}

		args[key] = sec
	}

	return nil
}
