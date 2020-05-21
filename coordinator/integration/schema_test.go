package integration

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestSchema(t *testing.T) {
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

	dbURL, err := primitives.NewDBURL("postgres://postgres:dev@my.cool.website:5432/mydb")
	gm.Expect(err).To(gm.BeNil())

	l, err := primitives.NewLabel("my-transactions")
	gm.Expect(err).To(gm.BeNil())

	source, err := client.AddSource(ctx, l, dbURL, nil)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create a new source", func(t *testing.T) {
		blob := map[string]interface{}{
			"my-transactions": map[string]interface{}{
				"col-1": "int",
				"col-2": "int",
				"col-3": "string",
			},
		}

		err = client.ReportSchema(ctx, source.ID, blob)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can update a source", func(t *testing.T) {
		blob := map[string]interface{}{
			"my-transactions": map[string]interface{}{
				"col-1": "int",
				"col-2": "int",
				"col-3": "string",
			},
		}

		err = client.ReportSchema(ctx, source.ID, blob)
		gm.Expect(err).To(gm.BeNil())

		// schema changed!
		blob = map[string]interface{}{
			"my-transactions": map[string]interface{}{
				"col-1": "int",
				"col-2": "int",
				"col-3": "string",
				"col-4": "string",
			},
		}

		err = client.ReportSchema(ctx, source.ID, blob)
		gm.Expect(err).To(gm.BeNil())
	})
}