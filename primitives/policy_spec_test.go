package primitives

import (
	"fmt"
	gm "github.com/onsi/gomega"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func loadPolicy() ([]byte, error) {
	path := filepath.Join("./testdata", "policy.yaml")
	return ioutil.ReadFile(path)
}

func TestYamlUnmarshalling(t *testing.T) {
	gm.RegisterTestingT(t)

	data, err := loadPolicy()
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

	where := rule.Where
	gm.Expect(len(where[0])).To(gm.Equal(1))
	gm.Expect(where[0]["partner"]).To(gm.Equal("visa"))

	field := rule.Fields[0]
	gm.Expect(field).To(gm.Equal(Field("card_number")))
}

func TestYamlMarshalling(t *testing.T) {
	gm.RegisterTestingT(t)

	spec := &PolicySpec{
		Version: PolicyVersion(1),
		Label:   "protect-ssn",
		Rules: []Rule{
			{
				Target: "records:creditcards.transactions",
				Action: Read,
				Effect: Deny,
				Fields: []Field{"card_number"},
				Where: []Where{
					{"partner": "visa"},
				},
			},
		},
	}

	d, err := spec.ToBytes()
	gm.Expect(err).To(gm.BeNil())

	expected, err := loadPolicy()
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(d).To(gm.Equal(expected), fmt.Sprintf("Wanted \n%s, got \n%s", string(expected), string(d)))
}
