package sources

import (
	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgtype"
)

// VariableSize represents the size for variable sized
// data types in postgres
const VariableSize = -1

var ci *pgtype.ConnInfo

// Data coming out of the postgres information schema
// is in the long form but pgx uses the short form.
// Set up any aliases required here.
var postgresDataTypeAliases = map[string]string{
	"timestamp without time zone": "timestamp",
	"smallint":                    "int2",
	"integer":                     "int4",
	"bigint":                      "int8",
	"real":                        "float4",
	"double precision":            "float8",
	"character":                   "bpchar",
	"character varying":           "varchar",
	"boolean":                     "bool",
}

// PostgresDataTypeToProtoField is a map to postgres OIDs
var PostgresDataTypeToProtoField = map[uint32]*proto.FieldInfo{
	// Note: serial2, serial4 and serial8 are not exposed outside of postgres
	// the types are just returned as int2, int4, or int8.
	pgtype.Int2OID: {
		Field: proto.FieldType_SMALLINT,
		Size:  2,
	},
	pgtype.Int4OID: {
		Field: proto.FieldType_INT,
		Size:  4,
	},
	pgtype.Int8OID: {
		Field: proto.FieldType_BIGINT,
		Size:  8,
	},
	pgtype.Float8OID: {
		Field: proto.FieldType_DOUBLE,
		Size:  8,
	},
	pgtype.Float4OID: {
		Field: proto.FieldType_REAL,
		Size:  4,
	},
	pgtype.BPCharOID: {
		Field: proto.FieldType_CHAR,
		Size:  VariableSize,
	},
	pgtype.VarcharOID: {
		Field: proto.FieldType_VARCHAR,
		Size:  VariableSize,
	},
	pgtype.TextOID: {
		Field: proto.FieldType_TEXT,
		Size:  VariableSize,
	},
	pgtype.TimestampOID: {
		Field: proto.FieldType_TIMESTAMP,
		Size:  8,
	},
	pgtype.BoolOID: {
		Field: proto.FieldType_BOOL,
		Size:  1,
	},
	pgtype.ByteaOID: {
		Field: proto.FieldType_BYTEA,
		Size:  VariableSize,
	},
}

type dataTypeToProtoFieldFunc func(dataType string) (*proto.FieldInfo, error)

var dataTypeToProtoFieldFuncs = map[primitives.SourceType]dataTypeToProtoFieldFunc{
	primitives.PostgresType: postgresDataTypeToProtoField,
}

// DataTypeToProtoField returns the field type from the postgres data type.
// Returns error if the type isn't in the map
func DataTypeToProtoField(sourceType primitives.SourceType, dataType string) (*proto.FieldInfo, error) {
	return dataTypeToProtoFieldFuncs[sourceType](dataType)
}

func postgresDataTypeToProtoField(dataType string) (*proto.FieldInfo, error) {
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
			Field: typ.Field,
		}
	}

	return fields, nil
}

func init() {
	ci = pgtype.NewConnInfo()
}
