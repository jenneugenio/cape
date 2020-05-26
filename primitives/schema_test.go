package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestSchema(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Validates a valid schema", func(t *testing.T) {
		sourceID, err := database.DecodeFromString("2015338ejcum4rzncvnugucvtc")
		gm.Expect(err).To(gm.BeNil())
		_, err = NewSchema(sourceID, SchemaBlob{
			"my-col": "INT",
		})

		gm.Expect(err).To(gm.BeNil())
	})
}
