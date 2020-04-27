package integration

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/coordinator/harness"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestTokens(t *testing.T) {
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

	t.Run("Create a token", func(t *testing.T) {
		gm.RegisterTestingT(t)

		token, err := client.NewToken(ctx, nil)
		gm.Expect(err).To(gm.BeNil())

		fmt.Println("make a token!", token)
	})
}
