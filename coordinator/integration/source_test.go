// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
)

func TestSource(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err).To(gm.BeNil())

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	workerToken, err := m.CreateWorker(ctx)
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
		gm.Expect(source.ServiceID).To(gm.BeNil())

		// Credentials are only returned for linked data connector
		gm.Expect(source.Credentials).To(gm.BeNil())
	})

	t.Run("create a new source with link", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("card-transactions")
		gm.Expect(err).To(gm.BeNil())

		emailStr := "service:connector@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		serviceURL, err := primitives.NewURL("https://localhost:8081")
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.DataConnectorServiceType, serviceURL)
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, &service.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(source.Service).ToNot(gm.BeNil())
		gm.Expect(source.Service.ID).To(gm.Equal(service.ID))
	})

	t.Run("change the link to an existing source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("new-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Type).To(gm.Equal(primitives.PostgresType))
		gm.Expect(source.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(source.ServiceID).To(gm.BeNil())
		gm.Expect(source.Credentials).To(gm.BeNil())

		emailStr := "service:another@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		serviceURL, err := primitives.NewURL("https://localhost:8081")
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.DataConnectorServiceType, serviceURL)
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(service.ID).ToNot(gm.BeNil())

		linkedSource, err := client.UpdateSource(ctx, source.Label, &service.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(linkedSource.Label).To(gm.Equal(l))
		gm.Expect(linkedSource.ID).ToNot(gm.BeNil())
		gm.Expect(linkedSource.Type).To(gm.Equal(primitives.PostgresType))
		gm.Expect(linkedSource.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(linkedSource.Service.ID).To(gm.Equal(service.ID))

		linkedSource, err = client.GetSource(ctx, source.ID, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(linkedSource.Label).To(gm.Equal(l))
		gm.Expect(linkedSource.ID).ToNot(gm.BeNil())
		gm.Expect(linkedSource.Type).To(gm.Equal(primitives.PostgresType))
		gm.Expect(linkedSource.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(linkedSource.Service.ID).To(gm.Equal(service.ID))
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

		service, err := primitives.NewService(email, primitives.UserServiceType, nil)
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

		out, err := client.GetSource(ctx, source.ID, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(out.Label).To(gm.Equal(label))
		gm.Expect(out.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(out.Credentials).To(gm.BeNil())
	})

	t.Run("pull a single data source by label", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("my-super-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
		gm.Expect(err).To(gm.BeNil())

		otherSource, err := client.GetSourceByLabel(ctx, l, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(otherSource.Label).To(gm.Equal(source.Label))
		gm.Expect(otherSource.Endpoint).To(gm.Equal(source.Endpoint))
		gm.Expect(otherSource.Credentials).To(gm.BeNil())
	})

	t.Run("describe a single data source", func(t *testing.T) {
		l, err := primitives.NewLabel("describe-me")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
		gm.Expect(err).To(gm.BeNil())

		// report a fake schema for this source
		err = m.ReportSchema(ctx, workerToken, source.ID, primitives.SchemaBlob{"my-table": {"my-col": "INT"}})
		gm.Expect(err).To(gm.BeNil())

		s, err := client.GetSource(ctx, source.ID, &coordinator.SourceOptions{WithSchema: true})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(s).ToNot(gm.BeNil())
		gm.Expect(s.Schema.Blob).To(gm.Equal(primitives.SchemaBlob{"my-table": {"my-col": "INT"}}))
	})

	t.Run("describe a single data source by label", func(t *testing.T) {
		l, err := primitives.NewLabel("describe-me-label")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, nil)
		gm.Expect(err).To(gm.BeNil())

		// report a fake schema for this source
		err = m.ReportSchema(ctx, workerToken, source.ID, primitives.SchemaBlob{"my-table": {"my-col": "INT"}})
		gm.Expect(err).To(gm.BeNil())

		s, err := client.GetSourceByLabel(ctx, l, &coordinator.SourceOptions{WithSchema: true})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(s).ToNot(gm.BeNil())
		gm.Expect(s.Schema.Blob).To(gm.Equal(primitives.SchemaBlob{"my-table": {"my-col": "INT"}}))
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

		_, err = client.GetSource(ctx, source.ID, nil)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("get credentials as linked data-connector", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("cool-source")
		gm.Expect(err).To(gm.BeNil())

		emailStr := "service:connector-cool@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		serviceURL, err := primitives.NewURL("https://localhost:8081")
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.DataConnectorServiceType, serviceURL)
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		apiToken, _, err := client.CreateToken(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL, &service.ID)
		gm.Expect(err).To(gm.BeNil())

		u, err := m.URL()
		gm.Expect(err).To(gm.BeNil())
		transport := coordinator.NewHTTPTransport(u, nil)
		serviceClient := coordinator.NewClient(transport)

		_, err = serviceClient.TokenLogin(ctx, apiToken)
		gm.Expect(err).To(gm.BeNil())

		source, err = serviceClient.GetSource(ctx, source.ID, nil)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Credentials.String()).To(gm.Equal(dbURL.String()))
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

	t.Run("List sources", func(t *testing.T) {

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
	})

	t.Run("List sources after token login", func(t *testing.T) {
		emailStr := "service:connector-cool@connector.com"
		email, err := primitives.NewEmail(emailStr)
		gm.Expect(err).To(gm.BeNil())

		serviceURL, err := primitives.NewURL("https://localhost:8081")
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.DataConnectorServiceType, serviceURL)
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		l1, err := primitives.NewLabel("my-cool-transactions")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddSource(ctx, l1, dbURL, &service.ID)
		gm.Expect(err).To(gm.BeNil())

		apiToken, _, err := client.CreateToken(ctx, m.Admin.User)
		gm.Expect(err).To(gm.BeNil())

		err = client.Logout(ctx, nil)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.TokenLogin(ctx, apiToken)
		gm.Expect(err).To(gm.BeNil())

		sources, err := client.ListSources(ctx)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(sources).ToNot(gm.BeNil())
	})
}
