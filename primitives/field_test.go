package primitives

import (
	"fmt"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestFieldNames(t *testing.T) {
	gm.RegisterTestingT(t)

	validNames := []string{
		"neat",
		"cool_table",
		"really_cool2",
		"wow_",
	}

	invalidNames := []string{
		"_not_good",
		"22222",
	}

	for _, n := range validNames {
		t.Run(fmt.Sprintf("Valid name: %s", n), func(t *testing.T) {
			gm.RegisterTestingT(t)
			_, err := NewField(n)
			gm.Expect(err).To(gm.BeNil())
		})
	}

	for _, n := range invalidNames {
		t.Run(fmt.Sprintf("Valid name: %s", n), func(t *testing.T) {
			gm.RegisterTestingT(t)
			_, err := NewField(n)
			gm.Expect(err.Error()).To(gm.Equal("invalid_field: field must start with a letter, and then only contain letters, numbers, or underscores"))
		})
	}
}
