package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
)

type SchemaBlob map[string]interface{}

func (s SchemaBlob) Validate() error {
	return nil
}

// Schema is the data schema for a data source
// It is linked with a source via the source label, the schema is currently stored as a json blob.
type Schema struct {
	*database.Primitive
	SourceID database.ID `json:"source_id"`
	Schema   SchemaBlob  `json:"source_schema"`
}

func (s *Schema) GetType() types.Type {
	return SchemaPrimitiveType
}

func (s *Schema) Validate() error {
	err := s.SourceID.Validate()
	if err != nil {
		return err
	}

	return s.Schema.Validate()
}

func NewSchema(sourceID database.ID, schema map[string]interface{}) (*Schema, error) {
	p, err := database.NewPrimitive(SchemaPrimitiveType)
	if err != nil {
		return nil, err
	}

	s := &Schema{
		Primitive: p,
		SourceID:  sourceID,
		Schema:    schema,
	}

	return s, s.Validate()
}
