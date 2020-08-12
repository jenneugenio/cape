package graph

import (
	"context"
	"github.com/capeprivacy/cape/auth"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
	"time"
)

// Fill in a context to work with the resolvers
// This function should only be used in tests
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
