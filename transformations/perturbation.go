package transformations

import (
	"math/rand"
	"time"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type PerturbationTransform struct {
	field        string
	min          float64
	max          float64
	seed         int64
	sourceSeeded rand.Source
	randInstance *rand.Rand
}

func (p *PerturbationTransform) perturbationFloat64(x float64) (float64, error) {
	noise := p.randInstance.Float64()*(p.max-p.min) + p.min
	y := x + noise

	return y, nil
}

func (p *PerturbationTransform) perturbationInt64(x int64) (int64, error) {
	noise := p.randInstance.Int63n(int64(p.max-p.min)) + int64(p.min)
	y := x + noise

	return y, nil
}

func (p *PerturbationTransform) Transform(schema *proto.Schema, input *proto.Record) error {
	field, err := GetField(schema, input, p.field)
	if err != nil {
		return err
	}

	output := &proto.Field{}
	switch t := field.GetValue().(type) {
	case *proto.Field_Int64:
		res, err := p.perturbationInt64(t.Int64)
		if err != nil {
			return err
		}
		output.Value = &proto.Field_Int64{Int64: res}
	case *proto.Field_Int32:
		res, err := p.perturbationInt64(int64(t.Int32))
		if err != nil {
			return err
		}
		output.Value = &proto.Field_Int32{Int32: int32(res)}

	case *proto.Field_Double:
		res, err := p.perturbationFloat64(t.Double)
		if err != nil {
			return err
		}
		output.Value = &proto.Field_Double{Double: res}
	case *proto.Field_Float:
		res, err := p.perturbationFloat64(float64(t.Float))
		if err != nil {
			return err
		}
		output.Value = &proto.Field_Float{Float: float32(res)}
	default:
		return errors.New(UnsupportedType, "Attempted to call %s transform on an unsupported type %T", p.Function(), t)
	}

	return SetField(schema, input, output, p.field)
}

func (p *PerturbationTransform) getMinMax(args Args) (float64, float64, error) {
	min, ok, err := args.LookupFloat64("min")
	if err != nil {
		return 0, 0, err
	}
	if !ok {
		return 0, 0, errors.New(MissingArgument, "Perturbation transformation expects a min argument")
	}

	max, ok, err := args.LookupFloat64("max")
	if err != nil {
		return 0, 0, err
	}
	if !ok {
		return 0, 0, errors.New(MissingArgument, "Perturbation transformation expects a max argument")
	}

	return min, max, nil
}

func (p *PerturbationTransform) Initialize(args Args) error {
	min, max, err := p.getMinMax(args)
	if err != nil {
		return err
	}
	p.min = min
	p.max = max

	if p.min > p.max {
		return errors.New(WrongArgument, "Min should be less than Max")
	}

	p.seed = time.Now().UnixNano()
	p.sourceSeeded = rand.NewSource(p.seed)
	p.randInstance = rand.New(p.sourceSeeded)

	return nil
}

func (p *PerturbationTransform) Validate(args Args) error {
	min, max, err := p.getMinMax(args)
	if err != nil {
		return err
	}

	if min > max {
		return errors.New(WrongArgument, "Min should be less than Max")
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
	return &PerturbationTransform{
		field: field,
		min:   0,
		max:   1,
	}, nil
}
