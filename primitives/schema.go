package primitives

import (
	"encoding/json"
	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/mitchellh/mapstructure"
	"io"
)

type SchemaBlob map[string]map[string]string

func (s *SchemaBlob) UnmarshalGQL(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(t, s); err != nil {
			return err
		}

		return s.Validate()
	default:
		return errors.New(InvalidPolicySpecCause, "Unable to unmarshal gql schema blob")
	}
}

func (s SchemaBlob) MarshalGQL(w io.Writer) {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	w.Write(j) //nolint: errcheck
}

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
	Blob     SchemaBlob  `json:"blob"`
}

func (s *Schema) GetType() types.Type {
	return SchemaPrimitiveType
}

func (s *Schema) Validate() error {
	err := s.SourceID.Validate()
	if err != nil {
		return err
	}

	return s.Blob.Validate()
}

func NewSchema(sourceID database.ID, schema SchemaBlob) (*Schema, error) {
	p, err := database.NewPrimitive(SchemaPrimitiveType)
	if err != nil {
		return nil, err
	}

	s := &Schema{
		Primitive: p,
		SourceID:  sourceID,
		Blob:      schema,
	}

	return s, s.Validate()
}
