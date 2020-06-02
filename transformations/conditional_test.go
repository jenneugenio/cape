package transformations

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestConditionals(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("test evaluate", func(t *testing.T) {
		c, err := NewConditional("hey == 0")
		gm.Expect(err).To(gm.BeNil())

		params := map[string]interface{}{"hey": 0}
		shouldFilter, err := c.Evaluate(params)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(shouldFilter).To(gm.BeTrue())
	})

	t.Run("get vars", func(t *testing.T) {
		c, err := NewConditional("hey == 0")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(c.Vars()).To(gm.Equal([]string{"hey"}))
	})

	t.Run("new error", func(t *testing.T) {
		_, err := NewConditional("\"hey == '0'\"")
		gm.Expect(err).NotTo(gm.BeNil())
	})

	errorTests := []struct {
		name       string
		expression string
		params     map[string]interface{}
	}{
		{
			name:       "no params",
			expression: "hey == 0",
			params:     map[string]interface{}{},
		},
		{
			name:       "wrong return type",
			expression: "hey + 5",
			params:     map[string]interface{}{"hey": 10},
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			c, err := NewConditional(test.expression)
			gm.Expect(err).To(gm.BeNil())

			_, err = c.Evaluate(test.params)
			gm.Expect(err).ToNot(gm.BeNil())
		})
	}
}
