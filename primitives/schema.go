package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
)

// QuerySchema is the data schema for a data source
// It is linked with a source via the source label, the schema is currently stored
// TODO -- make this comment relevant
// as a json blob. Since gql doesn't support maps, we may not be able to store this as jsonb???
type Schema struct {
	*database.Primitive
	SourceLabel  Label `json:"source_label"`
	Schema map[string]interface{} `json:"source_schema"`
}

func (s *Schema) GetType() types.Type {
	return SchemaPrimitiveType
}

func (s *Schema) Validate() error {
	return s.SourceLabel.Validate()
}

// TODO -- does this need to be encryptable?

func NewSchema(sourceLabel Label, schema map[string]interface{}) (*Schema, error) {
	p, err := database.NewPrimitive(SchemaPrimitiveType)
	if err != nil {
		return nil, err
	}

	s := &Schema{
		Primitive: p,
		SourceLabel: sourceLabel,
		Schema: schema,
	}

	return s, s.Validate()
}