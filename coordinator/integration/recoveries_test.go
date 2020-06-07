// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
)

func TestRecoveries(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err).To(gm.BeNil())

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can recover account successfully", func(t *testing.T) {

	})

	t.Run("can request account recovery for unknown email", func(t *testing.T) {

	})

	t.Run("can't recover an unknown email", func(t *testing.T) {

	})

	t.Run("can't recover with wrong id", func(t *testing.T) {

	})

	t.Run("can't recover account with wrong secret", func(t *testing.T) {

	})

	t.Run("non-worker can't list recoveries", func(t *testing.T) {

	})

	t.Run("non-worker can't delete recoveries", func(t *testing.T) {

	})

	t.Run("a worker can list recoveries", func(t *testing.T) {

	})

	t.Run("a worker can delete recoveries", func(t *testing.T) {

	})
}
