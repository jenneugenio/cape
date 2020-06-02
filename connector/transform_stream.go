package connector

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/connector/sources"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/capeprivacy/cape/transformations"
)

// TransformStream handles the transformations needed to be
// done on the flowing out of the connector.
type TransformStream struct {
	sources.Stream
	schema       *proto.Schema
	transforms   []transformations.Transformation
	conditionals []*transformations.Conditional
}

// NewTransformStream constructs the transforms and creates a stream.
func NewTransformStream(stream sources.Stream, schema *proto.Schema,
	transforms []*primitives.Transformation) (*TransformStream, error) {
	initTransforms := make([]transformations.Transformation, len(transforms))
	conditionals := make([]*transformations.Conditional, len(transforms))

	for i, t := range transforms {
		ctor := transformations.Get(t.Function)
		initT, err := ctor(t.Field.String())
		if err != nil {
			return nil, err
		}

		info, err := fieldToFieldInfo(t.Field.String(), schema)
		if err != nil {
			return nil, err
		}

		err = validateSupportedTypes(info.Field, initT)
		if err != nil {
			return nil, err
		}

		if t.Where != "" {
			c, err := transformations.NewConditional(t.Where.String())
			if err != nil {
				return nil, err
			}
			conditionals[i] = c
		}

		err = initT.Initialize(t.Args)
		if err != nil {
			return nil, err
		}

		initTransforms[i] = initT
	}

	return &TransformStream{
		Stream:       stream,
		schema:       schema,
		transforms:   initTransforms,
		conditionals: conditionals,
	}, nil
}

// Send handles doing the transformations and then forwarding the
// results records to grpc.
func (t *TransformStream) Send(record *proto.Record) error {
	if record.Schema != nil {
		t.schema = record.Schema
	}

	for i, transform := range t.transforms {
		index, err := fieldToFieldIndex(transform.Field(), t.schema)
		if err != nil {
			return err
		}

		if t.conditionals[i] != nil {
			vars := t.conditionals[i].Vars()

			params, err := fillParams(t.schema, record, vars)
			if err != nil {
				return err
			}

			shouldNotTransform, err := t.conditionals[i].Evaluate(params)
			if err != nil {
				return err
			}

			if shouldNotTransform {
				continue
			}
		}

		output, err := transform.Transform(record.Fields[index])
		if err != nil {
			return err
		}

		record.Fields[index] = output
	}
	return t.Stream.Send(record)
}

// validateSupportedTypes check to see if the given type is in the list of supported types
func validateSupportedTypes(typ proto.FieldType, tform transformations.Transformation) error {
	for _, supportedType := range tform.SupportedTypes() {
		if supportedType == typ {
			return nil
		}
	}

	return errors.New(transformations.UnsupportedType, "Attempted to call %s transform "+
		"on an unsupported type %s", tform.Function(), typ)
}

// fieldToFieldIndex returns the index of the field given the string
func fieldToFieldIndex(field string, schema *proto.Schema) (int, error) {
	for i, info := range schema.Fields {
		if field == info.Name {
			return i, nil
		}
	}

	return -1, errors.New(FieldNotFound, "Could not find field %s for target %s", field, schema.Target)
}

// fieldToFieldInfo returns the field info for the given field name
func fieldToFieldInfo(field string, schema *proto.Schema) (*proto.FieldInfo, error) {
	for _, info := range schema.Fields {
		if field == info.Name {
			return info, nil
		}
	}

	return nil, errors.New(FieldNotFound, "Could not find field %s for target %s", field, schema.Target)
}

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
			return nil, errors.New(FieldNotFound, "Could not evaluate where clause because '%s' is not a field in %s",
				v, schema.GetDataSource())
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
