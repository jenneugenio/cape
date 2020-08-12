package graph

import (
	"context"
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
	gm "github.com/onsi/gomega"
	"testing"
	"time"
)

type testDatabase struct {
	tokensDB tokensDB
}

func (t testDatabase) Roles() db.RoleDB               { panic("implement me") }
func (t testDatabase) Users() db.UserDB               { panic("implement me") }
func (t testDatabase) Projects() db.ProjectsDB        { panic("implement me") }
func (t testDatabase) Contributors() db.ContributorDB { panic("implement me") }
func (t testDatabase) Config() db.ConfigDB            { panic("implement me") }
func (t testDatabase) Secrets() db.SecretDB           { panic("implement me") }
func (t testDatabase) Session() db.SessionDB          { panic("implement me") }
func (t testDatabase) Recoveries() db.RecoveryDB      { panic("implement me") }

func (t testDatabase) Tokens() db.TokensDB { return &t.tokensDB }

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

func resolverContext(ctx context.Context, opts *ctxOptions) context.Context {
	if opts == nil {
		opts = &ctxOptions{}
	}

	opts.FillDefaults()

	// put a logger on the ctx
	l := fw.TestLogger()
	ctx = context.WithValue(ctx, fw.LoggerContextKey, *l)

	session := &models.Session{
		ID:        models.NewID(),
		UserID:    opts.userID,
		OwnerID:   models.NewID(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Token:     nil,
	}

	authSession := auth.Session{
		User: &models.User{
			ID:    opts.userID,
			Email: "user@cape.com",
			Name:  "User McPerson",
		},
		Session: session,
		Roles: models.UserRoles{
			Global: models.Role{
				ID:    opts.role.String(),
				Label: opts.role,
			},
		},
	}

	// put the session on the ctx
	ctx = context.WithValue(ctx, fw.SessionContextKey, &authSession)

	return ctx
}

type ctxOptions struct {
	role   models.Label
	userID string
}

func (c *ctxOptions) FillDefaults() {
	if c.role.String() == "" {
		c.role = models.AdminRole
	}

	if c.userID == "" {
		c.userID = "admin"
	}
}

func TestTokensCreate(t *testing.T) {
	gm.RegisterTestingT(t)

	resolver := &Resolver{
		Database: testDatabase{
			tokensDB: tokensDB{},
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
