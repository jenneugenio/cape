package transformations

import (
	"time"

	"github.com/Knetic/govaluate"
	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/golang/protobuf/ptypes"
)

type Condition string

func (c Condition) String() string {
	return string(c)
}

func (c Condition) Validate() error {
	if c.String() != "" {
		_, err := govaluate.NewEvaluableExpression(c.String())
		return err
	}

	return nil
}

type ConditionalTransformation struct {
	exp *govaluate.EvaluableExpression
	t   Transformation
}

func NewConditionalTransformation(expression Condition, t Transformation) (*ConditionalTransformation, error) {
	exp, err := govaluate.NewEvaluableExpression(expression.String())
	if err != nil {
		return nil, err
	}

	return &ConditionalTransformation{
		exp: exp,
		t:   t,
	}, nil
}

func (c ConditionalTransformation) Transform(schema *proto.Schema, record *proto.Record) error {
	params, err := fillParams(schema, record, c.exp.Vars())
	if err != nil {
		return err
	}

	res, err := c.exp.Evaluate(params)
	if err != nil {
		return err
	}

	shouldTransform, ok := res.(bool)
	if !ok {
		return errors.New(EvaluateBoolOnly, "Conditional expressions should only evaluate to booleans")
	}

	if shouldTransform {
		return c.t.Transform(schema, record)
	}

	return nil
}

func (c ConditionalTransformation) Initialize(args Args) error        { return c.t.Initialize(args) }
func (c ConditionalTransformation) Validate(args Args) error          { return c.t.Validate(args) }
func (c ConditionalTransformation) SupportedType() []proto.FieldType  { return c.t.SupportedTypes() }
func (c ConditionalTransformation) Function() string                  { return c.t.Function() }
func (c ConditionalTransformation) Field() string                     { return c.t.Field() }
func (c ConditionalTransformation) SupportedTypes() []proto.FieldType { return c.t.SupportedTypes() }

func fillParams(schema *proto.Schema, record *proto.Record, vars []string) (map[string]interface{}, error) {
	params := make(map[string]interface{}, len(vars))
	for _, v := range vars {
		found := false
		for i, field := range schema.Fields {
			if v == field.Name {
				i, err := fieldToInterface(record.GetFields()[i])
				if err != nil {
					return nil, err
				}
				params[v] = i
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New(FieldNotFound,
				"Could not evaluate where clause because '%s' is not a field in %s", v, schema.GetDataSource())
		}
	}

	return params, nil
}

func fieldToInterface(field *proto.Field) (interface{}, error) {
	switch t := field.GetValue().(type) {
	case *proto.Field_Bool:
		return t.Bool, nil
	case *proto.Field_Bytes:
		return t.Bytes, nil
	case *proto.Field_Double:
		return t.Double, nil
	case *proto.Field_Float:
		return t.Float, nil
	case *proto.Field_Int32:
		return t.Int32, nil
	case *proto.Field_Int64:
		return t.Int64, nil
	case *proto.Field_String_:
		return t.String_, nil
	case *proto.Field_Timestamp:
		tim, err := ptypes.Timestamp(t.Timestamp)
		if err != nil {
			return nil, err
		}
		return tim.Format(time.RFC3339), nil
	}

	return nil, errors.New(InvalidFieldType, "Invalid field type got %t", field.GetValue())
}
