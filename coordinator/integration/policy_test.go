// +build integration

package integration

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/capeprivacy/cape/models"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
)

func mockPolicy() (*models.Policy, error) {
	mockPolicy, err := ioutil.ReadFile("./testdata/policy.yaml")
	if err != nil {
		return nil, err
	}

	return models.ParsePolicy(mockPolicy)
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

	policy, err := mockPolicy()
	gm.Expect(err).To(gm.BeNil())

	t.Run("create policy", func(t *testing.T) {
		label := models.Label("hehedata")

		policy.Label = label
		policy, err := client.CreatePolicy(ctx, policy)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(policy.Label).To(gm.Equal(label))

		otherPolicy, err := client.GetPolicyByLabel(ctx, string(policy.Label))
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(policy.Label).To(gm.Equal(otherPolicy.Label))
		gm.Expect(policy.ID).To(gm.Equal(otherPolicy.ID))
	})

	t.Run("delete policy", func(t *testing.T) {
		label := models.Label("ds-dl-data")

		policy.Label = label

		newPolicy, err := client.CreatePolicy(ctx, policy)
		gm.Expect(err).To(gm.BeNil())

		err = client.DeletePolicy(ctx, string(label))
		gm.Expect(err).To(gm.BeNil())

		otherPolicy, err := client.GetPolicy(ctx, newPolicy.ID)
		gm.Expect(err).NotTo(gm.BeNil())
		gm.Expect(otherPolicy).To(gm.BeNil())
	})

	t.Run("get policy by label", func(t *testing.T) {
		l := "wow-policy-is-cool"
		label := models.Label(l)

		policy.Label = label

		_, err = client.CreatePolicy(ctx, policy)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.GetPolicyByLabel(ctx, l)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("cannot create the same policy twice", func(t *testing.T) {
		// label := models.Label("make-me-twice")

		// p1 := models.NewPolicy(label, spec)

		// _, err = client.CreatePolicy(ctx, &p1)
		// gm.Expect(err).To(gm.BeNil())

		// p2 := models.NewPolicy(label, spec)

		// _, err = client.CreatePolicy(ctx, &p2)
		// gm.Expect(err).ToNot(gm.BeNil())
		// // TODO(thor): This test was missing the .To(...) clause and seemed to be
		// // working but was a no-op. The returned error is losing the cause.
		// //gm.Expect(errors.CausedBy(err, database.DuplicateCause)).To(gm.BeTrue())
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

	policy, err := mockPolicy()
	gm.Expect(err).To(gm.BeNil())

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	labelStrs := []string{"cool-policy", "sad-policy", "wow-policy"}
	policies := make([]*models.Policy, 0, len(labelStrs))
	for _, labelStr := range labelStrs {
		label := models.Label(labelStr)

		policy.Label = label

		policy, err := client.CreatePolicy(ctx, policy)
		gm.Expect(err).To(gm.BeNil())

		policies = append(policies, policy)
	}

	otherPolicies, err := client.ListPolicies(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(otherPolicies).To(gm.Equal(policies))
}
