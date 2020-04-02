package sources

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

var testSourceType primitives.SourceType = "test"

type testSource struct {
	source *primitives.Source
}

func (t *testSource) Label() primitives.Label {
	return primitives.Label("test")
}
func (t *testSource) Type() primitives.SourceType {
	return testSourceType
}
func (t *testSource) Schema(_ context.Context, _ Query) (*proto.Schema, error) {
	return &proto.Schema{}, nil
}
func (t *testSource) Query(_ context.Context, _ Query, _ *proto.Schema, _ Stream) error {
	return nil
}
func (t *testSource) Close(_ context.Context) error {
	return nil
}

func newTestSource(s *primitives.Source) (Source, error) {
	return &testSource{
		source: s,
	}, nil
}

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
