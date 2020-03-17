// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/machinebox/graphql"
	gm "github.com/onsi/gomega"
)

func TestSessions(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client := graphql.NewClient("http://localhost:8081/v1/query")

	t.Run("create new login session", func(t *testing.T) {
		gm.RegisterTestingT(t)

		req := graphql.NewRequest(fmt.Sprintf(`
			mutation CreateLoginSession{
				createLoginSession(input: { email: "%s" }) {
					id
					identity_id
					expires_at
					type
					token
					credentials {
						salt
						alg
					}
				}
			}
		`, tc.User.Email))

		var resp2 struct {
			Session primitives.Session `json:"createLoginSession"`
		}

		err = client.Run(ctx, req, &resp2)
		gm.Expect(err).To(gm.BeNil())

		session := resp2.Session
		gm.Expect(session.IdentityID).To(gm.Equal(tc.User.ID))
		gm.Expect(session.ExpiresAt).To(gm.BeTemporally(">", time.Now().UTC()))

		// give this a extra 1 minute to compare against
		gm.Expect(session.ExpiresAt).To(gm.BeTemporally("<", time.Now().UTC().Add(time.Minute*6)))

		gm.Expect(session.Type).To(gm.Equal(primitives.Login))
		gm.Expect(session.Credentials.Salt).To(gm.Equal(tc.Credentials.Salt))
		gm.Expect(session.Token).ToNot(gm.BeNil())
	})

	t.Run("create new authentication session", func(t *testing.T) {
		gm.RegisterTestingT(t)

		req := graphql.NewRequest(fmt.Sprintf(`
			mutation CreateLoginSession{
				createLoginSession(input: { email: "%s" }) {
					token
					credentials {
						salt
					}
				}
			}
		`, tc.User.Email))

		var resp struct {
			Session primitives.Session `json:"createLoginSession"`
		}
		err = client.Run(ctx, req, &resp)
		gm.Expect(err).To(gm.BeNil())

		creds := tc.Credentials
		creds.Salt = resp.Session.Credentials.Salt

		sig, err := creds.Sign(resp.Session.Token)
		gm.Expect(err).To(gm.BeNil())

		req = graphql.NewRequest(fmt.Sprintf(`
			mutation CreateAuthSession{
				createAuthSession(input: { signature: "%s" }) {
					id
					identity_id
					expires_at
					type
					token
				}
			}
		`, sig.String()))

		req.Header.Add("Authorization", "Bearer "+resp.Session.Token.String())

		var resp2 struct {
			Session primitives.Session `json:"createAuthSession"`
		}

		err = client.Run(ctx, req, &resp2)
		gm.Expect(err).To(gm.BeNil())

		session := resp2.Session

		gm.Expect(session.IdentityID).To(gm.Equal(tc.User.ID))
		gm.Expect(session.ExpiresAt).To(gm.BeTemporally(">", time.Now().UTC()))

		// give this a extra 1 minute to compare against
		gm.Expect(session.ExpiresAt).To(gm.BeTemporally("<", time.Now().UTC().Add(time.Hour*24+time.Minute)))

		gm.Expect(session.Type).To(gm.Equal(primitives.Authenticated))
		gm.Expect(session.Token).ToNot(gm.BeNil())
	})

	t.Run("test fake user fails", func(t *testing.T) {
		gm.RegisterTestingT(t)

		req := graphql.NewRequest(fmt.Sprintf(`
			mutation CreateLoginSession{
				createLoginSession(input: { email: "%s" }) {
					token
					credentials {
						salt
					}
				}
			}
		`, "faker@fakerson.fake"))

		var resp struct {
			Session primitives.Session `json:"createLoginSession"`
		}
		err = client.Run(ctx, req, &resp)
		gm.Expect(err).To(gm.BeNil())

		creds := tc.Credentials
		creds.Salt = resp.Session.Credentials.Salt

		sig, err := creds.Sign(resp.Session.Token)
		gm.Expect(err).To(gm.BeNil())

		req = graphql.NewRequest(fmt.Sprintf(`
			mutation CreateAuthSession{
				createAuthSession(input: { signature: "%s" }) {
					id
				}
			}
		`, sig.String()))

		req.Header.Add("Authorization", "Bearer "+resp.Session.Token.String())

		var resp2 struct {
			Session primitives.Session `json:"createAuthSession"`
		}

		err = client.Run(ctx, req, &resp2)
		gm.Expect(err.Error()).To(gm.Equal("graphql: authentication_failure: Failed to authenticate"))
	})

	t.Run("test incorrect credentials", func(t *testing.T) {
		gm.RegisterTestingT(t)

		req := graphql.NewRequest(fmt.Sprintf(`
			mutation CreateLoginSession{
				createLoginSession(input: { email: "%s" }) {
					token
					credentials {
						salt
					}
				}
			}
		`, "faker@fakerson.fake"))

		var resp struct {
			Session primitives.Session `json:"createLoginSession"`
		}
		err = client.Run(ctx, req, &resp)
		gm.Expect(err).To(gm.BeNil())

		// get incorrect creds
		creds, err := auth.NewCredentials([]byte("my-devious-secret"), nil)
		gm.Expect(err).To(gm.BeNil())

		creds.Salt = resp.Session.Credentials.Salt

		sig, err := creds.Sign(resp.Session.Token)
		gm.Expect(err).To(gm.BeNil())

		req = graphql.NewRequest(fmt.Sprintf(`
			mutation CreateAuthSession{
				createAuthSession(input: { signature: "%s" }) {
					id
				}
			}
		`, sig.String()))

		req.Header.Add("Authorization", "Bearer "+resp.Session.Token.String())

		var resp2 struct {
			Session primitives.Session `json:"createAuthSession"`
		}

		err = client.Run(ctx, req, &resp2)
		gm.Expect(err.Error()).To(gm.Equal("graphql: authentication_failure: Failed to authenticate"))
	})
}
