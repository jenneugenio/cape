package transformations

import (
	"github.com/capeprivacy/cape/connector/proto"
)

type identity struct {
	field string
}

func (i *identity) Transform(input *proto.Field) (*proto.Field, error) {
	return input, nil
}

func (i *identity) Initialize(args Args) error {
	return nil
}

func (i *identity) Validate(args Args) error {
	return nil
}

func (i *identity) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_BIGINT,
		proto.FieldType_BOOL,
		proto.FieldType_BYTEA,
		proto.FieldType_INT,
		proto.FieldType_SMALLINT,
		proto.FieldType_DOUBLE,
		proto.FieldType_REAL,
		proto.FieldType_TEXT,
		proto.FieldType_TIMESTAMP,
		proto.FieldType_VARCHAR,
		proto.FieldType_CHAR,
	}
}

func (i *identity) Function() string {
	return "identity"
}

func (i *identity) Field() string {
	return i.field
}

func NewIdentityTransform(field string) (Transformation, error) {
	i := &identity{field: field}
	return i, nil
}
