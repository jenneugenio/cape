// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/machinebox/graphql"
	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
)

func TestUsers(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
	gm.Expect(err).To(gm.BeNil())

	client := graphql.NewClient("http://localhost:8081/v1/query")
	req := graphql.NewRequest(fmt.Sprintf(`
		mutation CreateUser {
		  createUser(input: { name: "Jerry Berry", email: "jerry@jerry.berry", public_key: "%s", salt: "%s", alg: "EDDSA"}) {
			id
			name
			email
		  }
		}
	`, creds.PublicKey.String(), creds.Salt.String()))

	var resp struct {
		User primitives.User `json:"createUser"`
	}

	err = client.Run(ctx, req, &resp)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(resp.User.Name).To(gm.Equal("Jerry Berry"))
	gm.Expect(resp.User.Email).To(gm.Equal("jerry@jerry.berry"))
}
