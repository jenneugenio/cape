package primitives

import (
	"fmt"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestTarget(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Valid target", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := NewTarget("records:collection.transactions")
		gm.Expect(err).To(gm.BeNil())
	})

	invalids := []string{
		"hello",
		"wow:cool",
		"this.shouldnt.work",
		"invalidtype:hmm.okay",
	}

	for _, invalid := range invalids {
		t.Run(fmt.Sprintf("Invalid target: %s", invalid), func(t *testing.T) {
			gm.RegisterTestingT(t)
			_, err := NewTarget(invalid)
			gm.Expect(err.Error()).To(gm.Equal("invalid_target: Target must be in the form <type>:<collection>.<collection>"))
		})
	}
}
