package transformations

import (
	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type plusOne struct {
	field string
}

func (p *plusOne) Transform(input *proto.Field) (*proto.Field, error) {
	output := &proto.Field{}
	switch t := input.GetValue().(type) {
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
		return nil, errors.New(UnsupportedType, "Attempted to call %s transform on an unsupported type %T", p.Function(), t)
	}

	return output, nil
}

func (p *plusOne) Initialize(args Args) error {
	return nil
}

func (p *plusOne) Validate(args Args) error {
	return nil
}

func (p *plusOne) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_BIGINT,
		proto.FieldType_INT,
		proto.FieldType_SMALLINT,
		proto.FieldType_DOUBLE,
		proto.FieldType_REAL,
	}
}

func (p *plusOne) Function() string {
	return "plusOne"
}

func (p *plusOne) Field() string {
	return p.field
}

func NewPlusOneTransform(field string) (Transformation, error) {
	p := &plusOne{field: field}
	return p, nil
}
