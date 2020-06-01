package transformations

import (
	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"math"
)

type RoundingType int

// See more at https://en.wikipedia.org/wiki/Rounding
// and https://www.cockroachlabs.com/blog/rounding-implementations-in-go/
const (
	RoundToEven RoundingType = iota
	RoundAwayFromZero
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
	return x, errors.New(UnsupportedType, "Unexpected error (unsupported rounding type %d)", r.roundingType)
}

func (r *RoundingTransform) Transform(input *proto.Field) (*proto.Field, error) {
	switch ty := input.GetValue().(type) {
	case *proto.Field_Double:
		res, err := r.roundFloat64(ty.Double)
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Double{Double: res}
		return output, nil
	case *proto.Field_Float:
		res, err := r.roundFloat64(float64(ty.Float))
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Float{Float: float32(res)}
		return output, nil
	}
	return input, nil
}

func (r *RoundingTransform) Initialize(args Args) error {
	roundingType, found := args["roundingType"]
	if found {
		switch roundingType.(string) {
		case "roundToEven":
			r.roundingType = RoundToEven
		case "awayFromZero":
			r.roundingType = RoundAwayFromZero
		default:
			return errors.New(UnsupportedType, "Unsupported rounding type '%s'", roundingType)
		}
	}

	precision, found := args["precision"]
	if found {
		r.precisionFactor = math.Pow10(precision.(int))
	}

	return nil
}

func (r *RoundingTransform) Validate(args Args) error {
	roundingType, found := args["roundingType"]
	if found {
		switch roundingType.(string) {
		case "roundToEven":
			break
		case "awayFromZero":
			break
		default:
			return errors.New(UnsupportedType, "Unsupported rounding type '%s'", roundingType)
		}
	}

	precision, found := args["precision"]
	if found {
		pre, ok := precision.(int)
		if !ok || pre < 0 {
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
