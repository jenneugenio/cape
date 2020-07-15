package auth

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/markbates/pkger"
	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database/types"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func TestDefaultAdminRBAC(t *testing.T) {
	gm.RegisterTestingT(t)

	policy, err := loadRBACFile(string(models.DefaultAdminRBAC) + ".yaml")
	gm.Expect(err).To(gm.BeNil())

	user := &models.User{}
	session, err := NewSession(user, &primitives.Session{}, []*models.RBACPolicy{policy},
		[]*primitives.Role{}, user)

	gm.Expect(err).To(gm.BeNil())
	gm.Expect(session).ToNot(gm.BeNil())

	testCases := []string{"users", "assignments", "attachments", "roles", "policies", "tokens"}

	for _, primitive := range testCases {
		t.Run(fmt.Sprintf("allowed create for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(models.Create, typ)
			gm.Expect(err).To(gm.BeNil())
		})

		t.Run(fmt.Sprintf("allowed delete for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(models.Delete, typ)
			gm.Expect(err).To(gm.BeNil())
		})
	}
}

func TestDefaultGlobalRBAC(t *testing.T) {
	gm.RegisterTestingT(t)

	policy, err := loadRBACFile(string(models.DefaultGlobalRBAC) + ".yaml")
	gm.Expect(err).To(gm.BeNil())

	user := &models.User{}
	session, err := NewSession(user, &primitives.Session{}, []*models.RBACPolicy{policy},
		[]*primitives.Role{}, user)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(session).ToNot(gm.BeNil())

	type TestCase struct {
		Primitive string
		Action    models.RBACAction
	}

	allowedTestCases := []TestCase{
		{"tokens", models.Create},
		{"tokens", models.Delete},
		{"sessions", models.Create},
		{"sessions", models.Delete},
		{"users", models.Read},
		{"attachments", models.Read},
		{"roles", models.Read},
		{"policies", models.Read},
	}

	for _, tc := range allowedTestCases {
		t.Run(fmt.Sprintf("allowed %s %s", tc.Action, tc.Primitive), func(t *testing.T) {
			typ, ok := types.Get(tc.Primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(tc.Action, typ)
			gm.Expect(err).To(gm.BeNil())
		})
	}

	deniedTestCases := []string{"users", "assignments", "attachments", "roles", "policies"}

	for _, primitive := range deniedTestCases {
		t.Run(fmt.Sprintf("denied create for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(models.Create, typ)
			gm.Expect(err).NotTo(gm.BeNil())
		})

		t.Run(fmt.Sprintf("denied delete for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(models.Delete, typ)
			gm.Expect(err).NotTo(gm.BeNil())
		})
	}
}

func loadRBACFile(file string) (*models.RBACPolicy, error) {
	dir := pkger.Dir("/primitives/policies/default")
	f, err := dir.Open(file)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return models.ParseRBACPolicy(b)
}
