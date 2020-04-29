package primitives

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func loadPolicy(file string) ([]byte, error) {
	path := filepath.Join("./testdata", file)
	return ioutil.ReadFile(path)
}

func TestYamlUnmarshalling(t *testing.T) {
	gm.RegisterTestingT(t)

	data, err := loadPolicy("policy.yaml")
	gm.Expect(err).To(gm.BeNil())

	spec, err := ParsePolicySpec(data)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(spec).ToNot(gm.BeNil())
	gm.Expect(spec.Version).To(gm.Equal(PolicyVersion(1)))
	gm.Expect(spec.Label.String()).To(gm.Equal("protect-ssn"))
	gm.Expect(len(spec.Rules)).To(gm.Equal(1))

	rule := spec.Rules[0]

	gm.Expect(rule.Target).To(gm.Equal(Target("records:creditcards.transactions")))
	gm.Expect(rule.Action).To(gm.Equal(Read))
	gm.Expect(rule.Effect).To(gm.Equal(Deny))
	gm.Expect(len(rule.Fields)).To(gm.Equal(1))

	field := rule.Fields[0]
	gm.Expect(field).To(gm.Equal(Field("card_number")))
}

func TestYamlMarshalling(t *testing.T) {
	gm.RegisterTestingT(t)

	spec := &PolicySpec{
		Version: PolicyVersion(1),
		Label:   "protect-ssn",
		Rules: []*Rule{
			{
				Target: "records:creditcards.transactions",
				Action: Read,
				Effect: Deny,
				Fields: []Field{"card_number"},
				Sudo:   false,
			},
		},
	}

	d, err := spec.ToBytes()
	gm.Expect(err).To(gm.BeNil())

	expected, err := loadPolicy("policy.yaml")
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(d).To(gm.Equal(expected), fmt.Sprintf("Wanted \n%s, got \n%s", string(expected), string(d)))
}

func TestPolicySpecRuleSudo(t *testing.T) {
	t.Run("policy with no sudo defaults false", func(t *testing.T) {
		gm.RegisterTestingT(t)
		data, err := loadPolicy("policy.yaml")
		gm.Expect(err).To(gm.BeNil())

		spec, err := ParsePolicySpec(data)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(spec.Rules[0].Sudo).To(gm.BeFalse())
	})

	t.Run("policy with sudo", func(t *testing.T) {
		gm.RegisterTestingT(t)
		data, err := loadPolicy("policy_with_sudo.yaml")
		gm.Expect(err).To(gm.BeNil())

		spec, err := ParsePolicySpec(data)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(spec.Rules[0].Sudo).To(gm.BeTrue())
	})
}

func TestPolicySpecValidation(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("cannot specify a where & fields clause in the same rule for records", func(t *testing.T) {
		gm.RegisterTestingT(t)
		spec := &PolicySpec{
			Version: PolicyVersion(1),
			Label:   "protect-ssn",
			Rules: []*Rule{
				{
					Target: "records:creditcards.transactions",
					Action: Read,
					Effect: Deny,
					Fields: []Field{"card_number"},
					Where: []Where{
						{"partner": "visa"},
					},
					Sudo: false,
				},
			},
		}

		err := spec.Validate()
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("invalid_policy_spec: Fields & Where cannot be specified on the same rule"))
	})

	t.Run("cannot specify a where or fields clause for internal policies", func(t *testing.T) {
		gm.RegisterTestingT(t)
		spec := &PolicySpec{
			Version: PolicyVersion(1),
			Label:   "protect-users",
			Rules: []*Rule{
				{
					Target: "internal:users.*",
					Action: Read,
					Effect: Deny,
					Fields: []Field{"FIELD"},
					Where: []Where{
						{"HEY": "HELLO"},
					},
				},
			},
		}

		err := spec.Validate()
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidPolicySpecCause)).To(gm.BeTrue())
	})
}

func TestPolicySpecGQL(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Test GQL Marshal", func(t *testing.T) {
		data, err := loadPolicy("policy_with_sudo.yaml")
		gm.Expect(err).To(gm.BeNil())

		spec, err := ParsePolicySpec(data)
		gm.Expect(err).To(gm.BeNil())

		by, err := json.Marshal(spec)
		gm.Expect(err).To(gm.BeNil())

		buf := &bytes.Buffer{}
		spec.MarshalGQL(buf)

		gm.Expect(buf.Bytes()).To(gm.Equal(by))
	})

	t.Run("Test GQL Unmarshal", func(t *testing.T) {
		data, err := loadPolicy("policy_with_sudo.yaml")
		gm.Expect(err).To(gm.BeNil())

		spec, err := ParsePolicySpec(data)
		gm.Expect(err).To(gm.BeNil())

		by, err := json.Marshal(spec)
		gm.Expect(err).To(gm.BeNil())

		var v map[string]interface{}
		err = json.Unmarshal(by, &v)
		gm.Expect(err).To(gm.BeNil())

		newSpec := PolicySpec{}
		err = newSpec.UnmarshalGQL(v)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(newSpec).To(gm.Equal(*spec))
	})
}
