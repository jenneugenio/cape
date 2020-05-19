package connector

import (
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
	schema     *proto.Schema
	transforms []transformations.Transformation
}

// NewTransformStream constructs the transforms and creates a stream.
func NewTransformStream(stream sources.Stream, schema *proto.Schema,
	transforms []*primitives.Transformation) (*TransformStream, error) {
	initTransforms := make([]transformations.Transformation, len(transforms))

	for i, t := range transforms {
		ctor := transformations.Get(t.Function.String())
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

		err = initT.Initialize(t.Args)
		if err != nil {
			return nil, err
		}

		initTransforms[i] = initT
	}

	return &TransformStream{
		Stream:     stream,
		schema:     schema,
		transforms: initTransforms,
	}, nil
}

// Send handles doing the transformations and then forwarding the
// results records to grpc.
func (t *TransformStream) Send(record *proto.Record) error {
	if record.Schema != nil {
		t.schema = record.Schema
	}

	for _, transform := range t.transforms {
		index, err := fieldToFieldIndex(transform.Field(), t.schema)
		if err != nil {
			return err
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
