package primitives

import (
	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestDisplayName(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Cannot have a short name", func(t *testing.T) {
		_, err := NewDisplayName("x")
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectNameCause)).To(gm.BeTrue())
	})

	t.Run("Cannot have a long name", func(t *testing.T) {
		name := ""
		for i := 0; i < 100; i++ {
			name += "b"
		}

		_, err := NewDisplayName(name)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectNameCause)).To(gm.BeTrue())
	})

	t.Run("Valid Name", func(t *testing.T) {
		d, err := NewDisplayName("Cool Project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(d).To(gm.Equal(DisplayName("Cool Project")))
	})
}
