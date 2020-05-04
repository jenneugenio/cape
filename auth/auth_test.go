package auth

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/markbates/pkger"
	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database/types"
	"github.com/capeprivacy/cape/primitives"
)

func TestDefaultAdminPolicy(t *testing.T) {
	gm.RegisterTestingT(t)

	policy, err := loadPolicyFile(primitives.DefaultAdminPolicy.String() + ".yaml")
	gm.Expect(err).To(gm.BeNil())

	session, err := NewSession(&primitives.User{}, &primitives.Session{}, []*primitives.Policy{policy},
		[]*primitives.Role{})
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(session).ToNot(gm.BeNil())

	testCases := []string{"users", "services", "assignments", "attachments", "roles", "policies", "sources", "tokens"}

	for _, primitive := range testCases {
		t.Run(fmt.Sprintf("allowed create for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(primitives.Create, typ)
			gm.Expect(err).To(gm.BeNil())
		})

		t.Run(fmt.Sprintf("allowed delete for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(primitives.Delete, typ)
			gm.Expect(err).To(gm.BeNil())
		})
	}
}

func TestDefaultGlobalPolicy(t *testing.T) {
	gm.RegisterTestingT(t)

	policy, err := loadPolicyFile(primitives.DefaultGlobalPolicy.String() + ".yaml")
	gm.Expect(err).To(gm.BeNil())

	session, err := NewSession(&primitives.User{}, &primitives.Session{}, []*primitives.Policy{policy},
		[]*primitives.Role{})
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(session).ToNot(gm.BeNil())

	type TestCase struct {
		Primitive string
		Action    primitives.Action
	}

	allowedTestCases := []TestCase{
		{"tokens", primitives.Create},
		{"tokens", primitives.Delete},
		{"sessions", primitives.Create},
		{"sessions", primitives.Delete},
		{"users", primitives.Read},
		{"services", primitives.Read},
		{"attachments", primitives.Read},
		{"roles", primitives.Read},
		{"policies", primitives.Read},
		{"sources", primitives.Read},
	}

	for _, tc := range allowedTestCases {
		t.Run(fmt.Sprintf("allowed %s %s", tc.Action, tc.Primitive), func(t *testing.T) {
			typ, ok := types.Get(tc.Primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(tc.Action, typ)
			gm.Expect(err).To(gm.BeNil())
		})
	}

	deniedTestCases := []string{"users", "services", "assignments", "attachments", "roles", "policies", "sources"}

	for _, primitive := range deniedTestCases {
		t.Run(fmt.Sprintf("denied create for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(primitives.Create, typ)
			gm.Expect(err).NotTo(gm.BeNil())
		})

		t.Run(fmt.Sprintf("denied delete for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(primitives.Delete, typ)
			gm.Expect(err).NotTo(gm.BeNil())
		})
	}
}

func TestDefaultDataConnectorPolicy(t *testing.T) {
	gm.RegisterTestingT(t)

	policy, err := loadPolicyFile(primitives.DefaultDataConnectorPolicy.String() + ".yaml")
	gm.Expect(err).To(gm.BeNil())

	session, err := NewSession(&primitives.User{}, &primitives.Session{}, []*primitives.Policy{policy},
		[]*primitives.Role{})
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(session).ToNot(gm.BeNil())

	t.Run("allowed to read sources", func(t *testing.T) {
		typ, ok := types.Get("sources")
		gm.Expect(ok).To(gm.BeTrue())

		err := session.Can(primitives.Read, typ)
		gm.Expect(err).To(gm.BeNil())
	})

	deniedTestCases := []string{"users", "services", "assignments", "attachments", "roles", "policies", "sources"}

	for _, primitive := range deniedTestCases {
		t.Run(fmt.Sprintf("denied create for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(primitives.Create, typ)
			gm.Expect(err).NotTo(gm.BeNil())
		})

		t.Run(fmt.Sprintf("denied delete for %s", primitive), func(t *testing.T) {
			typ, ok := types.Get(primitive)
			gm.Expect(ok).To(gm.BeTrue())

			err := session.Can(primitives.Delete, typ)
			gm.Expect(err).NotTo(gm.BeNil())
		})
	}
}

func loadPolicyFile(file string) (*primitives.Policy, error) {
	dir := pkger.Dir("/primitives/policies/default")
	f, err := dir.Open(file)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return primitives.ParsePolicy(b)
}