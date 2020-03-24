// +build integration

package integration

import (
	"context"
	"github.com/dropoutlabs/cape/database"
	"net/url"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/controller"
)

func TestSource(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck
	var id database.ID

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = client.Login(ctx, tc.User.Email, tc.UserPassword)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create a new source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, "my-transactions", u)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal("my-transactions"))
		gm.Expect(source.ID).ToNot(gm.BeNil())
		gm.Expect(source.Endpoint.String()).To(gm.Equal("postgres://my.cool.website.com:5432/mydb"))

		id = source.ID
	})

	t.Run("pull your data sources", func(t *testing.T) {
		gm.RegisterTestingT(t)

		sources, err := client.ListSources(ctx)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(sources)).To(gm.Equal(1))
		gm.Expect(sources[0].Label).To(gm.Equal("my-transactions"))
	})

	t.Run("pull a single data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		source, err := client.GetSource(ctx, id)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source.Label).To(gm.Equal("my-transactions"))
		gm.Expect(source.Endpoint.String()).To(gm.Equal("postgres://my.cool.website.com:5432/mydb"))
	})

	t.Run("insert the same data source", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, err := url.Parse("postgres://postgres:dev@my.cool.website.com:5432/mydb")
		gm.Expect(err).To(gm.BeNil())

		source, err := client.AddSource(ctx, "my-transactions", u)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(source).To(gm.BeNil())
	})
}
