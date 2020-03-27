package primitives

import (
	gm "github.com/onsi/gomega"
	"testing"
)

func TestName(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Cannot create a single letter name", func(t *testing.T) {
		_, err := NewName("A")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Cannot create a massive name", func(t *testing.T) {
		_, err := NewName("Adjasjdklsaj dja jkdl ajkl djsakl djklas jdklas jdkl asjdkas jdklas jdlks ajdklas jdklasjdsa jdkla ")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Can create a 2 character name", func(t *testing.T) {
		_, err := NewName("Al")
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can create a first and last name", func(t *testing.T) {
		_, err := NewName("Al McGoon")
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can have a middle name", func(t *testing.T) {
		_, err := NewName("Sandwich McLettuce Jr")
		gm.Expect(err).To(gm.BeNil())
	})
}