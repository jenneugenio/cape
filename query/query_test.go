package query

import (
	gm "github.com/onsi/gomega"
	"testing"
)

func TestQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	type TestCase struct {
		// test name
		name string

		// rule/policy info
		target string
		effect Effect
		fields []string
		where  map[string]string

		// input & expected query
		input    string
		expected string
	}

	testCases := []*TestCase{
		{
			"It redacts a field you cannot access",
			"transactions",
			Deny,
			[]string{"processor"},
			nil,
			"SELECT processor, card_number, value FROM transactions",
			"SELECT card_number, value FROM transactions",
		},

		{
			"It can give you access to only things you can have",
			"transactions",
			Allow,
			[]string{"processor"},
			nil,
			"SELECT processor, card_number, value FROM transactions",
			"SELECT processor FROM transactions",
		},

		// TODO -- we need schema information to do the inverse of this! (deny)
		{
			"It can rewrite a star command",
			"transactions",
			Allow,
			[]string{"processor", "card_number", "processor"},
			nil,
			"SELECT * FROM transactions",
			"SELECT processor, card_number, processor FROM transactions",
		},

		{
			"It can filter based on row",
			"transactions",
			Allow,
			[]string{"card_number"},
			map[string]string{
				"processor": "visa",
			},
			"SELECT * FROM transactions",
			"SELECT card_number FROM transactions WHERE processor = 'visa'",
		},

		{
			"It can filter multiple conditions",
			"transactions",
			Allow,
			[]string{"card_number"},
			map[string]string{
				"processor": "visa",
				"vendor":    "Cool Shirts Inc.",
			},
			"SELECT * FROM transactions",
			"SELECT card_number FROM transactions WHERE processor = 'visa' AND vendor = 'Cool Shirts Inc.'",
		},

		// TODO -- what about them writing conditions??
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gm.RegisterTestingT(t)

			r, err := NewRule(tc.target, tc.effect, tc.where, tc.fields...)
			gm.Expect(err).To(gm.BeNil())
			p, err := NewPolicy(r)
			gm.Expect(err).To(gm.BeNil())

			q, err := Parse(tc.input)
			gm.Expect(err).To(gm.BeNil())

			q, err = q.Rewrite(p)
			gm.Expect(err).To(gm.BeNil())

			gm.Expect(q.Raw()).To(gm.Equal(tc.expected))
		})
	}

	t.Run("Errors when you can't access anything you've asked for", func(t *testing.T) {
		gm.RegisterTestingT(t)

		r, err := NewRule("transactions", Allow, nil, "processor")
		gm.Expect(err).To(gm.BeNil())
		p, err := NewPolicy(r)
		gm.Expect(err).To(gm.BeNil())

		// I've only asked for things I can't see!
		q, err := Parse("SELECT card_number, value FROM transactions")
		gm.Expect(err).To(gm.BeNil())

		_, err = q.Rewrite(p)
		gm.Expect(err.Error()).To(gm.Equal("no_possible_fields: Cannot access any requested fields"))
	})
}
