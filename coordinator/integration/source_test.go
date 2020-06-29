// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

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

	dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website:5432/mydb")
	gm.Expect(err).To(gm.BeNil())

	endpoint := "postgres://my.cool.website:5432/mydb"

	t.Run("create a new source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal(l))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Type).To(gm.Equal(primitives.PostgresType))
		gm.Expect(source.Endpoint.String()).To(gm.Equal(endpoint))

		// Credentials are only returned for linked data connector
		gm.Expect(source.Credentials).To(gm.BeNil())
	})

	t.Run("pull a single data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("a-single-source")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, label, dbURL)
		gm.Expect(err).To(gm.BeNil())

		out, err := client.GetSource(ctx, source.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(out.Label).To(gm.Equal(label))
		gm.Expect(out.Endpoint.String()).To(gm.Equal(endpoint))
		gm.Expect(out.Credentials).To(gm.BeNil())
	})

	t.Run("pull a single data source by label", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("my-super-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL)
		gm.Expect(err).To(gm.BeNil())

		otherSource, err := client.GetSourceByLabel(ctx, l)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(otherSource.Label).To(gm.Equal(source.Label))
		gm.Expect(otherSource.Endpoint).To(gm.Equal(source.Endpoint))
		gm.Expect(otherSource.Credentials).To(gm.BeNil())
	})

	t.Run("describe a single data source", func(t *testing.T) {
		l, err := primitives.NewLabel("describe-me")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL)
		gm.Expect(err).To(gm.BeNil())

		s, err := client.GetSource(ctx, source.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(s).ToNot(gm.BeNil())
	})

	t.Run("describe a single data source by label", func(t *testing.T) {
		l, err := primitives.NewLabel("describe-me-label")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddSource(ctx, l, dbURL)
		gm.Expect(err).To(gm.BeNil())

		s, err := client.GetSourceByLabel(ctx, l)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(s).ToNot(gm.BeNil())
	})

	t.Run("insert the same data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(source).To(gm.BeNil())
	})

	t.Run("delete a source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		l, err := primitives.NewLabel("delete-me")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, l, dbURL)
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

	t.Run("List sources", func(t *testing.T) {

		dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		l1, err := primitives.NewLabel("my-transactions")
		gm.Expect(err).To(gm.BeNil())

		l2, err := primitives.NewLabel("my-other-transactions")
		gm.Expect(err).To(gm.BeNil())

		source1, err := client.AddSource(ctx, l1, dbURL)
		gm.Expect(err).To(gm.BeNil())

		source2, err := client.AddSource(ctx, l2, dbURL)
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

		service, err := primitives.NewService(email, primitives.UserServiceType)
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		l1, err := primitives.NewLabel("my-cool-transactions")
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AddSource(ctx, l1, dbURL)
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
