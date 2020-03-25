// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
)

func TestPolicies(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = client.Login(ctx, tc.User.Email, tc.UserPassword)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create policy", func(t *testing.T) {
		label := "admin-disallowed"
		p, err := primitives.NewPolicy(label)
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
		label := "ds-dl-data"
		p, err := primitives.NewPolicy(label)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		err = client.DeletePolicy(ctx, policy.ID)
		gm.Expect(err).To(gm.BeNil())

		otherPolicy, err := client.GetPolicy(ctx, policy.ID)
		gm.Expect(err).NotTo(gm.BeNil())
		gm.Expect(otherPolicy).To(gm.BeNil())
	})
}

func TestListPolicies(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = client.Login(ctx, tc.User.Email, tc.UserPassword)
	gm.Expect(err).To(gm.BeNil())

	labels := []string{"cool-policy", "sad-policy", "wow-policy"}
	policies := make([]*primitives.Policy, 3)
	for i, label := range labels {
		p, err := primitives.NewPolicy(label)
		gm.Expect(err).To(gm.BeNil())

		policy, err := client.CreatePolicy(ctx, p)
		gm.Expect(err).To(gm.BeNil())

		policies[i] = policy
	}

	otherPolicies, err := client.ListPolicies(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(otherPolicies).To(gm.ContainElements(policies))
}
