package transformations

import (
	"testing"

	gm "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"
)

func TestTransformationArgs(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		name      string
		transform string
		args      Args
	}{
		{
			name:      "test token args",
			transform: "tokenization",
			args: Args{
				"maxSize": 10,
			},
		},
		{
			name:      "test perturbation args",
			transform: "perturbation",
			args: Args{
				"min": 10,
				"max": 20,
			},
		},
		{
			name:      "test rounding args",
			transform: "rounding",
			args: Args{
				"roundingType": "awayFromZero",
				"precision":    10,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			by, err := yaml.Marshal(test.args)
			gm.Expect(err).To(gm.BeNil())

			var args Args
			err = yaml.Unmarshal(by, &args)
			gm.Expect(err).To(gm.BeNil())

			ctor := Get(test.transform)
			gm.Expect(err).To(gm.BeNil())

			transform, err := ctor("test")
			gm.Expect(err).To(gm.BeNil())

			err = transform.Validate(args)
			gm.Expect(err).To(gm.BeNil())
		})
	}

	errorTests := []struct {
		name      string
		transform string
		args      Args
	}{
		{
			name:      "max size is a string",
			transform: "tokenization",
			args: Args{
				"maxSize": "not a size",
			},
		},
		{
			name:      "min is a string",
			transform: "perturbation",
			args: Args{
				"min":  "string",
				"max":  20,
				"seed": 1234,
			},
		},
		{
			name:      "max is a string",
			transform: "perturbation",
			args: Args{
				"min": 10,
				"max": "string",
			},
		},
		{
			name:      "roundingType is a number",
			transform: "rounding",
			args: Args{
				"roundingType": 10,
				"precision":    10,
			},
		},
		{
			name:      "precision is a string",
			transform: "rounding",
			args: Args{
				"roundingType": "awayFromZero",
				"precision":    "string",
			},
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			by, err := yaml.Marshal(test.args)
			gm.Expect(err).To(gm.BeNil())

			var args Args
			err = yaml.Unmarshal(by, &args)
			gm.Expect(err).To(gm.BeNil())

			ctor := Get(test.transform)
			gm.Expect(err).To(gm.BeNil())

			transform, err := ctor("test")
			gm.Expect(err).To(gm.BeNil())

			err = transform.Validate(args)
			gm.Expect(err).ToNot(gm.BeNil())
		})
	}
}
