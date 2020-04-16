package sources

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestCache(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	r := &Registry{}
	cfg := &Config{}
	err := r.Register(testSourceType, newTestSource)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can get source", func(t *testing.T) {
		cache, err := NewCache(cfg, &testClient{}, r)
		gm.Expect(err).To(gm.BeNil())

		src, err := cache.Get(ctx, primitives.Label("test"))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(src.Label()).To(gm.Equal(primitives.Label("test")))
		gm.Expect(src.Type()).To(gm.Equal(testSourceType))
	})

	t.Run("get source fails if client returns false", func(t *testing.T) {
		cache, err := NewCache(cfg, &errClient{}, r)
		gm.Expect(err).To(gm.BeNil())

		_, err = cache.Get(ctx, primitives.Label("test"))
		gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
	})

	t.Run("get source returns existing source", func(t *testing.T) {
		cache, err := NewCache(cfg, &testClient{}, r)
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
