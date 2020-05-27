package transformations

import (
	"math/rand"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type PerturbationTransform struct {
	field string
	min   float64
	max   float64
	seed  int
}

func (p *PerturbationTransform) perturbationFloat64(x float64) (float64, error) {
	s := rand.NewSource(int64(p.seed))
	r := rand.New(s)

	noise := r.Float64()*(p.max-p.min) + p.min
	y := x + noise

	return y, nil
}

func (p *PerturbationTransform) Transform(input *proto.Field) (*proto.Field, error) {
	switch t := input.GetValue().(type) {
	case *proto.Field_Int64:
		res, err := p.perturbationFloat64(float64(t.Int64))
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Int64{Int64: int64(res)}
		return output, nil

	case *proto.Field_Int32:
		res, err := p.perturbationFloat64(float64(t.Int32))
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Int32{Int32: int32(res)}
		return output, nil

	case *proto.Field_Double:
		res, err := p.perturbationFloat64(t.Double)
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Double{Double: res}
		return output, nil

	case *proto.Field_Float:
		res, err := p.perturbationFloat64(float64(t.Float))
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Float{Float: float32(res)}
		return output, nil
	}

	return input, nil
}

func (p *PerturbationTransform) Initialize(args Args) error {
	min, found := args["min"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a min argument")
	}

	var ok bool
	p.min, ok = min.(float64)
	if !ok {
		return errors.New(UnsupportedType, "Unsupported min: must be a double type")
	}

	max, found := args["max"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a max argument")
	}

	p.max, ok = max.(float64)
	if !ok {
		return errors.New(UnsupportedType, "Unsupported max: must be a double type")
	}

	if p.min > p.max {
		return errors.New(WrongArgument, "Min should be less than Max")
	}

	seed, found := args["seed"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a seed argument")
	}

	p.seed, ok = seed.(int)
	if !ok {
		return errors.New(UnsupportedType, "Unsupported seed: must be an integer type")
	}

	return nil
}

func (p *PerturbationTransform) Validate(args Args) error {
	const unsupportedTypeMsg = "Unsupported type for %q: found %q expected %q"

	min, found := args["min"]
	if found {
		switch v := min.(type) {
		case float64:
			break
		default:
			return errors.New(UnsupportedType, unsupportedTypeMsg, "min", v, "double")
		}
	}

	max, found := args["max"]
	if found {
		switch v := max.(type) {
		case float64:
			break
		default:
			return errors.New(UnsupportedType, unsupportedTypeMsg, "max", v, "double")
		}
	}

	seed, found := args["seed"]
	if found {
		switch v := seed.(type) {
		case int:
			break
		default:
			return errors.New(UnsupportedType, unsupportedTypeMsg, "seed", v, "int")
		}
	}
	return nil
}

func (p *PerturbationTransform) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_DOUBLE,
		proto.FieldType_REAL,
	}
}

func (p *PerturbationTransform) Function() string {
	return "perturbation"
}

func (p *PerturbationTransform) Field() string {
	return p.field
}

func NewPerturbationTransform(field string) (Transformation, error) {
	p := &PerturbationTransform{field: field, min: 0, max: 1, seed: 1234}
	return p, nil
}
