// +build integration

package integration

import (
	"context"
	"net/url"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/controller/harness"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/primitives"
)

func TestSource(t *testing.T) {
	gm.RegisterTestingT(t)

	var id database.ID

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create a new source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		l, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, u)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Endpoint.String()).To(gm.Equal("postgres://my.cool.website.com:5432/mydb"))

		id = source.ID
	})

	t.Run("pull your data sources", func(t *testing.T) {
		gm.RegisterTestingT(t)

		sources, err := client.ListSources(ctx)
		gm.Expect(err).To(gm.BeNil())

		l, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(sources)).To(gm.Equal(1))
		gm.Expect(sources[0].Label).To(gm.Equal(l))
	})

	t.Run("pull a single data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		source, err := client.GetSource(ctx, id)
		gm.Expect(err).To(gm.BeNil())

		expectedLabel, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(expectedLabel))
		gm.Expect(source.Endpoint.String()).To(gm.Equal("postgres://my.cool.website.com:5432/mydb"))
	})

	t.Run("insert the same data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())
		l, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, u)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(source).To(gm.BeNil())
	})

	t.Run("delete a source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("delete-me")
		gm.Expect(err).To(gm.BeNil())

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/deleteme")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddSource(ctx, l, u)
		gm.Expect(err).To(gm.BeNil())

		err = client.RemoveSource(ctx, l)
		gm.Expect(err).To(gm.BeNil())

		// Now, there should be no sources!
		sources, err := client.ListSources(ctx)
		gm.Expect(err).To(gm.BeNil())

		// 1 left from above
		gm.Expect(len(sources)).To(gm.Equal(1))
	})
}
