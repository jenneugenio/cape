package transformations

import (
	"github.com/capeprivacy/cape/connector/proto"
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
	Transform(input *proto.Field) (*proto.Field, error)

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
func Get(function string) Constructor {
	return registry[function]
}

func init() {
	registry = Registry(make(map[string]Constructor))

	// Add additional transforms here
	registry.Add("identity", NewIdentityTransform)
	registry.Add("plusOne", NewPlusOneTransform)
	registry.Add("rounding", NewRoundingTransform)
	registry.Add("perturbation", NewPerturbationTransform)
	registry.Add("tokenization", NewTokenizationTransform)
}
