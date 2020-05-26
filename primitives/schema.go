package primitives

import (
	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type SchemaBlob map[string]map[string]string

func (s SchemaBlob) Validate() error {
	for _, tableName := range s {
		for _, v := range tableName {
			_, ok := proto.FieldType_value[v]
			if !ok {
				return errors.New(UnsupportedSchemaCause, "Column type %s is unsupported", v)
			}
		}
	}

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

func NewSchema(sourceID database.ID, schema SchemaBlob) (*Schema, error) {
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
