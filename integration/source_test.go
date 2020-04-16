// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/primitives"
)

func TestSource(t *testing.T) {
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

	dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website:5432/mydb")
	gm.Expect(err).To(gm.BeNil())

	endpoint := "postgres://my.cool.website:5432/mydb"

	t.Run("create a new source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Type).To(gm.Equal(primitives.PostgresType))
		gm.Expect(source.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(source.Credentials.String()).To(gm.Equal(dbURL.String()))
		gm.Expect(source.ServiceID).To(gm.BeNil())
	})

	t.Run("create a new source with link", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("card-transactions")
		gm.Expect(err).To(gm.BeNil())

		emailStr := "service:connector@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("random-password"), nil)
		gm.Expect(err).To(gm.BeNil())

		serviceURL, err := primitives.NewURL("https://localhost:8081")
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.DataConnectorServiceType, serviceURL,
			creds.Package())
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, &service.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(source.Credentials.String()).To(gm.Equal(dbURL.String()))
		gm.Expect(source.ServiceID).ToNot(gm.BeNil())
		gm.Expect(*source.ServiceID).To(gm.Equal(service.ID))
	})

	t.Run("can't create link to a non-existent data connector", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("card-transactions-service-dooesnt-exist")
		gm.Expect(err).To(gm.BeNil())

		serviceID, err := database.GenerateID(primitives.ServicePrimitiveType)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddSource(ctx, l, dbURL, &serviceID)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("can't create link that is not a data connector", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("card-transactions")
		gm.Expect(err).To(gm.BeNil())

		emailStr := "service:user@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("random-password"), nil)
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.UserServiceType, nil, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, &service.ID)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(source).To(gm.BeNil())
	})

	t.Run("pull a single data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("a-single-source")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, label, dbURL, nil)
		gm.Expect(err).To(gm.BeNil())

		out, err := client.GetSource(ctx, source.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(out.Label).To(gm.Equal(label))
		gm.Expect(out.Endpoint.String()).To(gm.Equal(endpoint))
	})

	t.Run("pull a single data source by label", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("my-super-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
		gm.Expect(err).To(gm.BeNil())

		otherSource, err := client.GetSourceByLabel(ctx, l)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(otherSource.Label).To(gm.Equal(source.Label))
		gm.Expect(otherSource.Endpoint).To(gm.Equal(source.Endpoint))
	})

	t.Run("insert the same data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(source).To(gm.BeNil())
	})

	t.Run("delete a source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("delete-me")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
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

	dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website.com:5432/mydb")
	gm.Expect(err).To(gm.BeNil())

	l1, err := primitives.NewLabel("my-transactions")
	gm.Expect(err).To(gm.BeNil())

	l2, err := primitives.NewLabel("my-other-transactions")
	gm.Expect(err).To(gm.BeNil())

	source1, err := client.AddSource(ctx, l1, dbURL, nil)
	gm.Expect(err).To(gm.BeNil())

	source2, err := client.AddSource(ctx, l2, dbURL, nil)
	gm.Expect(err).To(gm.BeNil())

	sources, err := client.ListSources(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(len(sources)).To(gm.Equal(2))
	gm.Expect(sources).To(gm.ContainElements(source1, source2))
}
