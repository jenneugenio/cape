package graph

import (
	"context"
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/models"
	gm "github.com/onsi/gomega"
	"testing"
)

type tokensDB struct {
	// you can set return token as a default token to return in this test
	returnToken models.Token

	// params captured by the create/delete/list calls to be used for asserts
	createdToken models.Token
	deletedToken string
	listByUserID string
}

func (t *tokensDB) Get(ctx context.Context, s string) (*models.Token, error) {
	return &t.returnToken, nil
}

func (t *tokensDB) Create(ctx context.Context, token models.Token) error {
	t.createdToken = token
	return nil
}

func (t *tokensDB) Delete(ctx context.Context, s string) error {
	t.deletedToken = s
	return nil
}

func (t *tokensDB) ListByUserID(ctx context.Context, s string) ([]models.Token, error) {
	t.listByUserID = s
	return nil, nil
}

func TestTokensCreate(t *testing.T) {
	gm.RegisterTestingT(t)

	resolver := &Resolver{
		Database: testDatabase{
			tokensDB: &tokensDB{},
		},
		CredentialProducer: &auth.SHA256Producer{},
	}

	mutationResolver := resolver.Mutation()
	t.Run("admin can create a token for themselves", func(t *testing.T) {
		req := model.CreateTokenRequest{
			UserID: "admin",
		}

		ctx := resolverContext(context.TODO(), nil)
		resp, err := mutationResolver.CreateToken(ctx, req)

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.Token).ToNot(gm.BeNil())
		gm.Expect(resp.Token.UserID).To(gm.Equal("admin"))
		gm.Expect(resp.Secret).ToNot(gm.Equal(""))
	})

	t.Run("admin can create a token for anyone", func(t *testing.T) {
		req := model.CreateTokenRequest{
			UserID: "nottheadmin",
		}

		ctx := resolverContext(context.TODO(), nil)
		resp, err := mutationResolver.CreateToken(ctx, req)

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.Token).ToNot(gm.BeNil())
		gm.Expect(resp.Token.UserID).To(gm.Equal("nottheadmin"))
		gm.Expect(resp.Secret).ToNot(gm.Equal(""))
	})

	t.Run("non admin can create a token for themselves", func(t *testing.T) {
		req := model.CreateTokenRequest{
			UserID: "justsomeuser",
		}

		ctx := resolverContext(context.TODO(), &ctxOptions{
			role: models.UserRole, userID: "justsomeuser",
		})
		resp, err := mutationResolver.CreateToken(ctx, req)

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.Token).ToNot(gm.BeNil())
		gm.Expect(resp.Token.UserID).To(gm.Equal("justsomeuser"))
		gm.Expect(resp.Secret).ToNot(gm.Equal(""))
	})

	t.Run("non admin cannot create a token for someone else", func(t *testing.T) {
		req := model.CreateTokenRequest{
			UserID: "myfriend",
		}

		ctx := resolverContext(context.TODO(), &ctxOptions{role: models.UserRole, userID: "justsomeuser"})
		_, err := mutationResolver.CreateToken(ctx, req)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
