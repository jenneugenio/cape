package transformations

import (
	"math/rand"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type PerturbationTransform struct {
	field        string
	min          float64
	max          float64
	seed         int64
	sourceSeeded rand.Source
}

func (p *PerturbationTransform) perturbationFloat64(x float64) (float64, error) {
	r := rand.New(p.sourceSeeded)

	noise := r.Float64()*(p.max-p.min) + p.min
	y := x + noise

	return y, nil
}

func (p *PerturbationTransform) perturbationInt64(x int64) (int64, error) {
	r := rand.New(p.sourceSeeded)

	noise := r.Int63n(int64(p.max-p.min)) + int64(p.min)
	y := x + noise

	return y, nil
}

func (p *PerturbationTransform) Transform(input *proto.Field) (*proto.Field, error) {
	switch t := input.GetValue().(type) {
	case *proto.Field_Int64:
		res, err := p.perturbationInt64(t.Int64)
		if err != nil {
			return nil, err
		}
		output := &proto.Field{}
		output.Value = &proto.Field_Int64{Int64: res}
		return output, nil

	case *proto.Field_Int32:
		res, err := p.perturbationInt64(int64(t.Int32))
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
	const unsupportedTypeMsg = "Unsupported type for %q: found %q expected %q"

	min, found := args["min"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a min argument")
	}

	var ok bool
	p.min, ok = min.(float64)
	if !ok {
		return errors.New(UnsupportedType, unsupportedTypeMsg, "min", min, "double")
	}

	max, found := args["max"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a max argument")
	}

	p.max, ok = max.(float64)
	if !ok {
		return errors.New(UnsupportedType, unsupportedTypeMsg, "max", max, "double")
	}

	if p.min > p.max {
		return errors.New(WrongArgument, "Min should be less than Max")
	}

	seed, found := args["seed"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a seed argument")
	}

	p.seed, ok = seed.(int64)
	if !ok {
		return errors.New(UnsupportedType, unsupportedTypeMsg, "seed", seed, "int")
	}

	p.sourceSeeded = rand.NewSource(p.seed)

	return nil
}

func (p *PerturbationTransform) Validate(args Args) error {
	const unsupportedTypeMsg = "Unsupported type for %q: found %q expected %q"

	min, found := args["min"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a min argument")
	}

	var ok bool
	min, ok = min.(float64)
	if !ok {
		return errors.New(UnsupportedType, unsupportedTypeMsg, "min", min, "double")
	}

	max, found := args["max"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a max argument")
	}

	max, ok = max.(float64)
	if !ok {
		return errors.New(UnsupportedType, unsupportedTypeMsg, "max", max, "double")
	}

	if p.min > p.max {
		return errors.New(WrongArgument, "Min should be less than Max")
	}

	seed, found := args["seed"]
	if !found {
		return errors.New(MissingArgument, "Perturbation transformation expects a seed argument")
	}

	seed, ok = seed.(int64)
	if !ok {
		return errors.New(UnsupportedType, unsupportedTypeMsg, "seed", seed, "int")
	}

	return nil
}

func (p *PerturbationTransform) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_BIGINT,
		proto.FieldType_DOUBLE,
		proto.FieldType_INT,
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
