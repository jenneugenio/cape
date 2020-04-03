package sources

import (
	"github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// VariableSize represents the size for variable sized
// data types in postgres
const VariableSize = -1

// Field represents a proto field and its pg data type size
type Field struct {
	FieldType proto.FieldType

	// number of bytes, -1 if variables sized,
	// see constant VariableSize above
	Size int64
}

var postgresDataTypeToProtoField = map[string]*Field{
	"integer": &Field{ // nolint: gofmt
		FieldType: proto.FieldType_INT,
		Size:      4,
	},
	"text": &Field{ // nolint: gofmt
		FieldType: proto.FieldType_TEXT,
		Size:      VariableSize,
	},
	"timestamp without time zone": &Field{ // nolint: gofmt
		FieldType: proto.FieldType_TIMESTAMP,
		Size:      8,
	},
	"bigint": &Field{ // nolint: gofmt
		FieldType: proto.FieldType_BIGINT,
		Size:      8,
	},
	"double precision": &Field{ // nolint: gofmt
		FieldType: proto.FieldType_DOUBLE,
		Size:      8,
	},
}

var sourceTypeDataTypeRegistry = map[primitives.SourceType]map[string]*Field{
	primitives.PostgresType: postgresDataTypeToProtoField,
}

// DataTypeToProtoField returns the field type from the postgres data type.
// Returns error if the type isn't in the map
func DataTypeToProtoField(sourceType primitives.SourceType, dataType string) (*Field, error) {
	fieldType, ok := sourceTypeDataTypeRegistry[sourceType][dataType]
	if !ok {
		err := errors.New(UnknownFieldType, "Cannot convert %s to a known field type for source type %s.", dataType, sourceType)
		return nil, err
	}

	return fieldType, nil
}
