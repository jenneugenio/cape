package transformations

import (
	"fmt"
	"github.com/capeprivacy/cape/connector/proto"
	"math"
)

type RoundingType int

// See more at https://en.wikipedia.org/wiki/Rounding
// and https://www.cockroachlabs.com/blog/rounding-implementations-in-go/
const (
	RoundToEven RoundingType = iota
	RoundAwayFromZero
)

type rounding struct {
	field           string
	roundingType    RoundingType
	precisionFactor float64
}

func (p *rounding) roundFloat64(x float64) (float64, error) {
	y := x * p.precisionFactor
	switch p.roundingType {
	case RoundAwayFromZero:
		return math.Round(y) / p.precisionFactor, nil
	case RoundToEven:
		return math.RoundToEven(y) / p.precisionFactor, nil
	}
	return x, fmt.Errorf("Unsupported rounding type %d", p.roundingType)
}

func (p *rounding) Transform(input *proto.Field) (*proto.Field, error) {
	switch t := input.GetValue().(type) {
	case *proto.Field_Double:
		res, err := p.roundFloat64(t.Double)
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Double{Double: res}
		return output, nil
	case *proto.Field_Float:
		res, err := p.roundFloat64(float64(t.Float))
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Float{Float: float32(res)}
		return output, nil
	}
	return input, nil
}

func (p *rounding) Initialize(args Args) error {
	roundingType, found := args["roundingType"]
	if found {
		switch roundingType.(string) {
		case "roundToEven":
			p.roundingType = RoundToEven
		case "awayFromZero":
			p.roundingType = RoundAwayFromZero
		default:
			return fmt.Errorf("Unsupported rounding type '%s'", roundingType)
		}
	}

	precision, found := args["precision"]
	if found {
		p.precisionFactor = math.Pow10(precision.(int))
	}

	return nil
}

func (p *rounding) Validate(args Args) error {
	roundingType, found := args["roundingType"]
	if found {
		switch roundingType.(string) {
		case "roundToEven":
			break
		case "awayFromZero":
			break
		default:
			return fmt.Errorf("Unsupported rounding type '%s'", roundingType)
		}
	}

	precision, found := args["precision"]
	if found {
		switch precision.(type) {
		case int:
			break
		default:
			return fmt.Errorf("Unsupported precision: must be positive integer")
		}

		pre := precision.(int)
		switch {
		case pre >= 0:
			break
		default:
			return fmt.Errorf("Unsupported precision: must be positive integer")
		}
	}

	return nil
}

func (p *rounding) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_DOUBLE,
		proto.FieldType_REAL,
	}
}

func (p *rounding) Function() string {
	return "rounding"
}

func (p *rounding) Field() string {
	return p.field
}

func NewRoundingTransform(field string) (Transformation, error) {
	p := &rounding{field: field, precisionFactor: 1.0}
	return p, nil
}
