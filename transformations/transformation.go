package transformations

import (
	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Args represents the arguments to be passed into
// a transformation
type Args map[string]interface{}

var (
	// Registry contains all of the transformations currently
	// registered with cape.
	registry Registry
)

// Transformation represents a transformation that will be
// called when required by the policy
type Transformation interface {
	// Transform does the transform work. Its up the implementation of
	// transform to handle the different input types and return an
	// error if an invalid type is encountered. The output type will
	// be handled by the grpc encoding.
	Transform(schema *proto.Schema, input *proto.Record) error

	// Initialize does any complex and potential expensive work to prepare
	// for transforming data.
	Initialize(args Args) error

	// Validate validates whether the transformation is valid given
	// the args passed in.
	Validate(args Args) error

	// SupportedTypes returns all the types that a transformation supports
	SupportedTypes() []proto.FieldType

	// Function returns the name of the function name of the transformation
	Function() string

	// Field returns the field the transformation is being applied to
	Field() string
}

// Constructor function to create a Transformation. Field is the field labe
// that the transformation will be applied to.
type Constructor func(field string) (Transformation, error)

// Registry is a registry of all constructors used to create
// different transformations.
type Registry map[string]Constructor

// Add adds a new transformation contructor to the registry
func (c Registry) Add(function string, ctor Constructor) {
	c[function] = ctor
}

// Get returns the Transformation constructor for the given function label
func Get(function string) (Constructor, error) {
	ctor, ok := registry[function]
	if !ok {
		return nil, errors.New(TransformationNotFound, "Could not find transformation %s", function)
	}

	return ctor, nil
}

func init() {
	registry = Registry(make(map[string]Constructor))

	// Add additional transforms here
	registry.Add("identity", NewIdentityTransform)
	registry.Add("plusOne", NewPlusOneTransform)
	registry.Add("rounding", NewRoundingTransform)
	registry.Add("perturbation", NewPerturbationTransform)
	registry.Add("tokenization", NewTokenizationTransform)
	registry.Add("scrambler", NewScramblerTransform)
}

func GetField(schema *proto.Schema, record *proto.Record, field string) (*proto.Field, error) {
	i, err := fieldToFieldIndex(schema, field)
	if err != nil {
		return nil, err
	}

	return record.Fields[i], nil
}

func SetField(schema *proto.Schema, record *proto.Record, newField *proto.Field, fieldName string) error {
	i, err := fieldToFieldIndex(schema, fieldName)
	if err != nil {
		return err
	}

	record.Fields[i] = newField
	return nil
}

// fieldToFieldIndex returns the index of the field given the string
func fieldToFieldIndex(schema *proto.Schema, field string) (int, error) {
	for i, info := range schema.Fields {
		if field == info.Name {
			return i, nil
		}
	}

	return -1, errors.New(FieldNotFound, "Could not find field %s for target %s", field, schema.Target)
}
