// +build integration

package integration

import (
	"context"
	"net/url"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
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

		source, err := client.AddSource(ctx, l, u, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Endpoint.String()).To(gm.Equal("postgres://my.cool.website.com:5432/mydb"))
		gm.Expect(source.ServiceID).To(gm.Equal(database.EmptyID))

		id = source.ID
	})

	t.Run("create a new source with link", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		l, err := primitives.NewLabel("card-transactions")
		gm.Expect(err).To(gm.BeNil())

		emailStr := "service:connector@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("random-password"), nil)
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.DataConnectorServiceType, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, u, &service.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Endpoint.String()).To(gm.Equal("postgres://my.cool.website.com:5432/mydb"))
		gm.Expect(source.ServiceID).To(gm.Equal(service.ID))
	})

	t.Run("can't create link that is not a data connector", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		l, err := primitives.NewLabel("card-transactions")
		gm.Expect(err).To(gm.BeNil())

		emailStr := "service:user@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("random-password"), nil)
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.UserServiceType, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, u, &service.ID)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(source).To(gm.BeNil())
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

		source, err := client.AddSource(ctx, l, u, nil)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(source).To(gm.BeNil())
	})

	t.Run("delete a source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("delete-me")
		gm.Expect(err).To(gm.BeNil())

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/deleteme")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, u, nil)
		gm.Expect(err).To(gm.BeNil())

		err = client.RemoveSource(ctx, l)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.GetSource(ctx, source.ID)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}

func TestListSources(t *testing.T) {
	gm.RegisterTestingT(t)

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

	u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/mydb")
	gm.Expect(err).To(gm.BeNil())

	l1, err := primitives.NewLabel("my-transactions")
	gm.Expect(err).To(gm.BeNil())

	l2, err := primitives.NewLabel("my-other-transactions")
	gm.Expect(err).To(gm.BeNil())

	source1, err := client.AddSource(ctx, l1, u, nil)
	gm.Expect(err).To(gm.BeNil())

	source2, err := client.AddSource(ctx, l2, u, nil)
	gm.Expect(err).To(gm.BeNil())

	sources, err := client.ListSources(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(len(sources)).To(gm.Equal(2))
	gm.Expect(sources).To(gm.ContainElements(source1, source2))
}
