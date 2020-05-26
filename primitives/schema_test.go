package primitives

import (
	"encoding/json"
	"github.com/capeprivacy/cape/coordinator/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestSchema(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Validates a valid schema", func(t *testing.T) {
		sourceID, err := database.DecodeFromString("2015338ejcum4rzncvnugucvtc")
		gm.Expect(err).To(gm.BeNil())
		_, err = NewSchema(sourceID, SchemaBlob{
			"my-table": {"my-col": "INT"},
		})

		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Cannot create an invalid schema", func(t *testing.T) {
		sourceID, err := database.DecodeFromString("2015338ejcum4rzncvnugucvtc")
		gm.Expect(err).To(gm.BeNil())
		_, err = NewSchema(sourceID, SchemaBlob{
			"my-table": {"my-col": "not a real data type!!! :O!111"},
		})

		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, UnsupportedSchemaCause)).To(gm.BeTrue())
	})

	t.Run("Can unmarshal from a string", func(t *testing.T) {
		blob := []byte(`{ "my-table" : { "my-col": "INT" }}`)
		var schemaBlob SchemaBlob
		err := json.Unmarshal(blob, &schemaBlob)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(schemaBlob["my-table"]["my-col"]).To(gm.Equal("INT"))
	})
}
