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
		blob := primitives.SchemaBlob{
			"my-transactions": {
				"col-1": "INT",
				"col-2": "INT",
				"col-3": "TEXT",
			},
		}

		err = client.ReportSchema(ctx, source.ID, blob)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can update a source", func(t *testing.T) {
		blob := primitives.SchemaBlob{
			"my-transactions": {
				"col-1": "INT",
				"col-2": "INT",
				"col-3": "TEXT",
			},
		}

		err = client.ReportSchema(ctx, source.ID, blob)
		gm.Expect(err).To(gm.BeNil())

		// schema changed!
		blob = primitives.SchemaBlob{
			"my-transactions": {
				"col-1": "INT",
				"col-2": "INT",
				"col-3": "TEXT",
				"col-4": "TEXT",
			},
		}

		err = client.ReportSchema(ctx, source.ID, blob)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Cannot report invalid column types", func(t *testing.T) {
		blob := primitives.SchemaBlob{
			"my-transactions": {
				"col-1": "thiskindofdatatypeisprobablynotfoundinmostdatabases",
			},
		}

		err = client.ReportSchema(ctx, source.ID, blob)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
