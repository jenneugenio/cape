package transformations

import (
	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type PlusOneTransform struct {
	field string
}

func (p *PlusOneTransform) Transform(schema *proto.Schema, input *proto.Record) error {
	field, err := GetField(schema, input, p.field)
	if err != nil {
		return err
	}

	output := &proto.Field{}
	switch t := field.GetValue().(type) {
	case *proto.Field_Double:
		res := t.Double + 1
		output.Value = &proto.Field_Double{Double: res}
	case *proto.Field_Float:
		res := t.Float + 1
		output.Value = &proto.Field_Float{Float: res}
	case *proto.Field_Int32:
		res := t.Int32 + 1
		output.Value = &proto.Field_Int32{Int32: res}
	case *proto.Field_Int64:
		res := t.Int64 + 1
		output.Value = &proto.Field_Int64{Int64: res}
	default:
		return errors.New(UnsupportedType, "Attempted to call %s transform on an unsupported type %T", p.Function(), t)
	}

	return SetField(schema, input, output, p.field)
}

func (p *PlusOneTransform) Initialize(args Args) error {
	return nil
}

func (p *PlusOneTransform) Validate(args Args) error {
	return nil
}

func (p *PlusOneTransform) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_BIGINT,
		proto.FieldType_INT,
		proto.FieldType_SMALLINT,
		proto.FieldType_DOUBLE,
		proto.FieldType_REAL,
	}
}

func (p *PlusOneTransform) Function() string {
	return "plusOne"
}

func (p *PlusOneTransform) Field() string {
	return p.field
}

func NewPlusOneTransform(field string) (Transformation, error) {
	return &PlusOneTransform{field: field}, nil
}
