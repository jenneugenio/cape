package integration

import (
	"context"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/machinebox/graphql"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestUsers(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	id, err := database.DecodeFromString("3m0v6dgh7avzj9gdm63qnp819x")
	gm.Expect(err).To(gm.BeNil())

	client := graphql.NewClient("http://localhost:8081/query")
	req := graphql.NewRequest(`
		mutation CreateUser {
		  createUser(input: { name: "Jerry", id: "3m0v6dgh7avzj9gdm63qnp819x" }) {
			id
			name
		  }
		}
	`)

	var resp struct {
		User primitives.User `json:"createUser"`
	}

	err = client.Run(ctx, req, &resp)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(resp.User.Name).To(gm.Equal("Jerry"))
	gm.Expect(resp.User.ID).To(gm.Equal(id))
}
