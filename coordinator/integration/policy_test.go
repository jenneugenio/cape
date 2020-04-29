// +build integration

package integration

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	"io/ioutil"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
)

func mockSpec() (*primitives.PolicySpec, error) {
	mockSpec, err := ioutil.ReadFile("./testdata/policy.yaml")
	if err != nil {
		return nil, err
	}

	return primitives.ParsePolicySpec(mockSpec)
}

func TestPolicies(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	spec, err := mockSpec()
	gm.Expect(err).To(gm.BeNil())

	t.Run("create policy", func(t *testing.T) {
		label, err := primitives.NewLabel("admin-disallowed")
		gm.Expect(err).To(gm.BeNil())

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(policy.Label).To(gm.Equal(label))

		otherPolicy, err := client.GetPolicy(ctx, policy.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(policy.Label).To(gm.Equal(otherPolicy.Label))
		gm.Expect(policy.ID).To(gm.Equal(otherPolicy.ID))
	})

	t.Run("delete policy", func(t *testing.T) {
		label, err := primitives.NewLabel("ds-dl-data")
		gm.Expect(err).To(gm.BeNil())

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		err = client.DeletePolicy(ctx, policy.ID)
		gm.Expect(err).To(gm.BeNil())

		otherPolicy, err := client.GetPolicy(ctx, policy.ID)
		gm.Expect(err).NotTo(gm.BeNil())
		gm.Expect(otherPolicy).To(gm.BeNil())
	})

	t.Run("get policy by label", func(t *testing.T) {
		label, err := primitives.NewLabel("wow-policy-is-cool")
		gm.Expect(err).To(gm.BeNil())

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.GetPolicyByLabel(ctx, label)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("cannot create the same policy twice", func(t *testing.T) {
		label, err := primitives.NewLabel("make-me-twice")
		gm.Expect(err).To(gm.BeNil())

		p1, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.CreatePolicy(ctx, p1)
		gm.Expect(err).To(gm.BeNil())

		p2, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.CreatePolicy(ctx, p2)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, database.DuplicateCause))
	})
}

func TestListPolicies(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	spec, err := mockSpec()
	gm.Expect(err).To(gm.BeNil())

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	labelStrs := []string{"cool-policy", "sad-policy", "wow-policy"}
	labels := make([]primitives.Label, len(labelStrs))

	policies := make([]*primitives.Policy, 3)
	for i, labelStr := range labelStrs {
		label, err := primitives.NewLabel(labelStr)
		gm.Expect(err).To(gm.BeNil())

		labels[i] = label

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		policies[i] = policy
	}

	otherPolicies, err := client.ListPolicies(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(otherPolicies).To(gm.ContainElements(policies))
}

func TestAttachments(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	spec, err := mockSpec()
	gm.Expect(err).To(gm.BeNil())

	t.Run("attach policy", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("admin-disallowed")
		gm.Expect(err).To(gm.BeNil())

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		role, err := client.CreateRole(ctx, "owner", nil)
		gm.Expect(err).To(gm.BeNil())

		attachment, err := client.AttachPolicy(ctx, policy.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(attachment.Policy.Label).To(gm.Equal(policy.Label))
		gm.Expect(attachment.Role.Label).To(gm.Equal(role.Label))

		policies, err := client.GetRolePolicies(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(policies)).To(gm.Equal(1))
	})

	t.Run("detach policy", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("ds-allowed")
		gm.Expect(err).To(gm.BeNil())

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		role, err := client.CreateRole(ctx, "data-scientist", nil)
		gm.Expect(err).To(gm.BeNil())

		attachment, err := client.AttachPolicy(ctx, policy.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())

		err = client.DetachPolicy(ctx, attachment.Policy.ID, attachment.Role.ID)
		gm.Expect(err).To(gm.BeNil())

		policies, err := client.GetRolePolicies(ctx, role.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(policies)).To(gm.Equal(0))
	})

	t.Run("test get policies for identity", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("cio-allowed")
		gm.Expect(err).To(gm.BeNil())

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		role, err := client.CreateRole(ctx, "cioo", nil)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AssignRole(ctx, m.Admin.User.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AttachPolicy(ctx, policy.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())

		policies, err := client.GetIdentityPolicies(ctx, m.Admin.User.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(policies)).To(gm.Equal(1))
		gm.Expect(policies[0].Label).To(gm.Equal(label))
	})

	t.Run("attach policy twice", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label, err := primitives.NewLabel("attachme")
		gm.Expect(err).To(gm.BeNil())

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		role, err := client.CreateRole(ctx, "coolguy", nil)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AttachPolicy(ctx, policy.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AttachPolicy(ctx, policy.ID, role.ID)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, database.DuplicateCause))
	})
}
