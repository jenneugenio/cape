package sources

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

type testClient struct{}

func (t *testClient) GetSource(ctx context.Context, source primitives.Label) (*primitives.Source, error) {
	return &primitives.Source{
		Label: source,
		Type:  testSourceType,
	}, nil
}

type errClient struct{}

func (e *errClient) GetSource(ctx context.Context, source primitives.Label) (*primitives.Source, error) {
	return nil, errors.New(NotFoundCause, "whoops")
}

func TestCache(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	r := &Registry{}
	err := r.Register(testSourceType, newTestSource)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can get source", func(t *testing.T) {
		cache, err := NewCache(&testClient{}, r)
		gm.Expect(err).To(gm.BeNil())

		src, err := cache.Get(ctx, primitives.Label("test"))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(src.Label()).To(gm.Equal(primitives.Label("test")))
		gm.Expect(src.Type()).To(gm.Equal(testSourceType))
	})

	t.Run("get source fails if client returns false", func(t *testing.T) {
		cache, err := NewCache(&errClient{}, r)
		gm.Expect(err).To(gm.BeNil())

		_, err = cache.Get(ctx, primitives.Label("test"))
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("get source returns existing source", func(t *testing.T) {
		cache, err := NewCache(&testClient{}, r)
		gm.Expect(err).To(gm.BeNil())

		src, err := cache.Get(ctx, primitives.Label("test"))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(src.Label()).To(gm.Equal(primitives.Label("test")))
		gm.Expect(src.Type()).To(gm.Equal(testSourceType))

		srcTwo, err := cache.Get(ctx, primitives.Label("test"))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(srcTwo).To(gm.Equal(src))
	})

	t.Run("close, closes all sources", func(t *testing.T) {

	})
}
