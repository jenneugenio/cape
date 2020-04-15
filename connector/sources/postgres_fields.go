package sources

import (
	"github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgtype"
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

var ci *pgtype.ConnInfo

// Data coming out of the postgre information schema
// is in the long form but pgx uses the short form.
// Set up any aliases required here.
var postgresDataTypeAliases = map[string]string{
	"timestamp without time zone": "timestamp",
	"bigint":                      "int8",
	"double precision":            "float8",
	"integer":                     "int4",
}

// PostgresDataTypeToProtoField is a map to postgres OIDs
var PostgresDataTypeToProtoField = map[uint32]*Field{
	pgtype.Int4OID: {
		FieldType: proto.FieldType_INT,
		Size:      4,
	},
	pgtype.TextOID: {
		FieldType: proto.FieldType_TEXT,
		Size:      VariableSize,
	},
	pgtype.TimestampOID: {
		FieldType: proto.FieldType_TIMESTAMP,
		Size:      8,
	},
	pgtype.Int8OID: {
		FieldType: proto.FieldType_BIGINT,
		Size:      8,
	},
	pgtype.Float8OID: {
		FieldType: proto.FieldType_DOUBLE,
		Size:      8,
	},
}

type dataTypeToProtoFieldFunc func(dataType string) (*Field, error)

var dataTypeToProtoFieldFuncs = map[primitives.SourceType]dataTypeToProtoFieldFunc{
	primitives.PostgresType: postgresDataTypeToProtoField,
}

// DataTypeToProtoField returns the field type from the postgres data type.
// Returns error if the type isn't in the map
func DataTypeToProtoField(sourceType primitives.SourceType, dataType string) (*Field, error) {
	return dataTypeToProtoFieldFuncs[sourceType](dataType)
}

func postgresDataTypeToProtoField(dataType string) (*Field, error) {
	var dType *pgtype.DataType

	// try an alias
	alias, ok := postgresDataTypeAliases[dataType]
	if ok {
		dataType = alias
	}

	dType, ok = ci.DataTypeForName(dataType)
	if !ok {
		err := errors.New(UnknownFieldType, "Cannot convert data type alias %s to a known field type for source type %s.",
			alias, primitives.PostgresType)
		return nil, err
	}

	fieldType, ok := PostgresDataTypeToProtoField[dType.OID]
	if !ok {
		err := errors.New(UnknownFieldType, "Cannot convert oid %d to a known field type for source type %s.",
			dType.OID, primitives.PostgresType)
		return nil, err
	}

	return fieldType, nil
}

func fieldsFromFieldDescription(fds []pgproto3.FieldDescription) ([]*proto.FieldInfo, error) {
	fields := make([]*proto.FieldInfo, len(fds))

	for i, fd := range fds {
		typ, ok := PostgresDataTypeToProtoField[fd.DataTypeOID]
		if !ok {
			err := errors.New(UnknownFieldType, "Cannot find type %d for postgres oid type", fd.DataTypeOID)
			return nil, err
		}

		fields[i] = &proto.FieldInfo{
			Name:  string(fd.Name),
			Size:  int64(fd.DataTypeSize),
			Field: typ.FieldType,
		}
	}

	return fields, nil
}

func init() {
	ci = pgtype.NewConnInfo()
}
