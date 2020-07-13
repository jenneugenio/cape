package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

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
	ID      string `json:"id"`
	Version uint8  `json:"version"`
	Label   Label  `json:"label"`

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
