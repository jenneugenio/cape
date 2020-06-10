package transformations

import (
	"math"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type RoundingType string

// See more at https://en.wikipedia.org/wiki/Rounding
// and https://www.cockroachlabs.com/blog/rounding-implementations-in-go/
const (
	RoundToEven       RoundingType = "roundToEven"
	RoundAwayFromZero RoundingType = "awayFromZero"
)

type RoundingTransform struct {
	field           string
	roundingType    RoundingType
	precisionFactor float64
}

func (r *RoundingTransform) roundFloat64(x float64) (float64, error) {
	y := x * r.precisionFactor
	switch r.roundingType {
	case RoundAwayFromZero:
		return math.Round(y) / r.precisionFactor, nil
	case RoundToEven:
		return math.RoundToEven(y) / r.precisionFactor, nil
	}
	return x, errors.New(UnsupportedType, "Unexpected error (unsupported rounding type %s)", r.roundingType)
}

func (r *RoundingTransform) Transform(schema *proto.Schema, input *proto.Record) error {
	field, err := GetField(schema, input, r.field)
	if err != nil {
		return err
	}

	output := &proto.Field{}
	switch ty := field.GetValue().(type) {
	case *proto.Field_Double:
		res, err := r.roundFloat64(ty.Double)
		if err != nil {
			return err
		}
		output.Value = &proto.Field_Double{Double: res}
	case *proto.Field_Float:
		res, err := r.roundFloat64(float64(ty.Float))
		if err != nil {
			return err
		}
		output.Value = &proto.Field_Float{Float: float32(res)}
	default:
		return errors.New(UnsupportedType, "Attempted to call %s transform on an unsupported type %T", r.Function(), ty)
	}

	return SetField(schema, input, output, r.field)
}

func (r *RoundingTransform) Initialize(args Args) error {
	roundingType, found := args["roundingType"]
	if found {
		rt, ok := roundingType.(string)
		if !ok {
			return errors.New(UnsupportedType, "Rounding args expected roundingType string saw %T", roundingType)
		}

		switch RoundingType(rt) {
		case RoundToEven:
			r.roundingType = RoundToEven
		case RoundAwayFromZero:
			r.roundingType = RoundAwayFromZero
		default:
			return errors.New(UnsupportedType, "Unsupported rounding type '%s'", roundingType)
		}
	}

	precision, found := args["precision"]
	if found {
		r.precisionFactor = math.Pow10(int(precision.(float64)))
	}

	return nil
}

func (r *RoundingTransform) Validate(args Args) error {
	roundingType, found := args["roundingType"]
	if found {
		rt, ok := roundingType.(string)
		if !ok {
			return errors.New(UnsupportedType, "Rounding args expected roundingType string saw %T", roundingType)
		}

		switch RoundingType(rt) {
		case RoundToEven:
			break
		case RoundAwayFromZero:
			break
		default:
			return errors.New(UnsupportedType, "Unsupported rounding type '%s'", roundingType)
		}
	}

	precision, found, err := args.LookupFloat64("precision")
	if err != nil {
		return err
	}
	if found {
		if precision < 0 {
			return errors.New(UnsupportedType, "Unsupported precision: must be positive integer")
		}
	}

	return nil
}

func (r *RoundingTransform) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_DOUBLE,
		proto.FieldType_REAL,
	}
}

func (r *RoundingTransform) Function() string {
	return "rounding"
}

func (r *RoundingTransform) Field() string {
	return r.field
}

func NewRoundingTransform(field string) (Transformation, error) {
	r := &RoundingTransform{
		field:           field,
		roundingType:    RoundToEven,
		precisionFactor: 1.0,
	}
	return r, nil
}
