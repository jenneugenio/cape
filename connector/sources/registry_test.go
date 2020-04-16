package sources

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestSourcesRegistry(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Can register a source", func(t *testing.T) {
		r := &Registry{}
		err := r.Register(testSourceType, newTestSource)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can get registered sources", func(t *testing.T) {
		r := &Registry{}
		err := r.Register(testSourceType, newTestSource)
		gm.Expect(err).To(gm.BeNil())

		ctor, err := r.Get(testSourceType)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(ctor).ToNot(gm.BeNil())
	})

	t.Run("Returns error for unknown source", func(t *testing.T) {
		r := &Registry{}
		_, err := r.Get(testSourceType)
		gm.Expect(errors.FromCause(err, SourceNotSupported)).To(gm.BeTrue())
	})
}
