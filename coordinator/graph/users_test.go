package graph

import (
	"context"
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/models"
	gm "github.com/onsi/gomega"
	"testing"
)

type usersDB struct {
	createdUser *models.User
}

func (u *usersDB) Create(ctx context.Context, user models.User) error {
	u.createdUser = &user
	return nil
}

func (u *usersDB) Update(ctx context.Context, s string, user models.User) error {
	panic("implement me")
}

func (u *usersDB) Delete(ctx context.Context, email models.Email) (db.DeleteStatus, error) {
	panic("implement me")
}

func (u *usersDB) Get(ctx context.Context, email models.Email) (*models.User, error) {
	panic("implement me")
}

func (u *usersDB) GetByID(ctx context.Context, s string) (*models.User, error) {
	panic("implement me")
}

func (u *usersDB) List(ctx context.Context, options *db.ListUserOptions) ([]models.User, error) {
	panic("implement me")
}

type rolesDB struct {
	orgRole models.Label
}

func (r *rolesDB) Get(ctx context.Context, label models.Label) (*models.Role, error) {
	panic("implement me")
}

func (r *rolesDB) GetByID(ctx context.Context, s string) (*models.Role, error) {
	panic("implement me")
}

func (r *rolesDB) List(ctx context.Context, options *db.ListRoleOptions) ([]*models.Role, error) {
	panic("implement me")
}

func (r *rolesDB) GetAll(ctx context.Context, s string) (*models.UserRoles, error) {
	panic("implement me")
}

func (r *rolesDB) SetOrgRole(ctx context.Context, email models.Email, label models.Label) (*models.Assignment, error) {
	r.orgRole = label
	return nil, nil
}

func (r *rolesDB) GetOrgRole(ctx context.Context, email models.Email) (*models.Role, error) {
	panic("implement me")
}

func (r *rolesDB) SetProjectRole(ctx context.Context, email models.Email, label models.Label, label2 models.Label) (*models.Assignment, error) {
	panic("implement me")
}

func (r *rolesDB) GetProjectRole(ctx context.Context, email models.Email, s string) (*models.Role, error) {
	panic("implement me")
}

func (r *rolesDB) CreateSystemRoles(ctx context.Context) error {
	panic("implement me")
}

func TestUsers(t *testing.T) {
	gm.RegisterTestingT(t)
	udb := &usersDB{}
	rdb := &rolesDB{}
	resolver := &Resolver{
		Database: testDatabase{
			usersDB: udb,
			rolesDB: rdb,
		},
		CredentialProducer: &auth.SHA256Producer{},
	}

	mutationResolver := resolver.Mutation()

	t.Run("admin can create a user", func(t *testing.T) {
		ctx := resolverContext(context.TODO(), nil)
		resp, err := mutationResolver.CreateUser(ctx, model.CreateUserRequest{
			Name:  "My Friend",
			Email: "coolguy@hotmail.com",
		})

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.User).ToNot(gm.BeNil())
		gm.Expect(udb.createdUser.Name.String()).To(gm.Equal("My Friend"))
		gm.Expect(rdb.orgRole).To(gm.Equal(models.UserRole))
	})
}
