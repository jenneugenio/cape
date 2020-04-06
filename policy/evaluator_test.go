package policy

import (
	"github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/dropoutlabs/cape/query"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestEvaluator(t *testing.T) {
	gm.RegisterTestingT(t)

	schema := &proto.Schema{
		DataSource: "transactions",
		Target:     "transactions",
		Type:       0,
		Fields: []*proto.Field{
			{Name: "id"},
			{Name: "processor"},
			{Name: "timestamp"},
			{Name: "card_id"},
			{Name: "card_number"},
			{Name: "value"},
			{Name: "ssn"},
			{Name: "vendor"},
		},
	}

	t.Run("Fails if no policies are attached", func(t *testing.T) {
		gm.RegisterTestingT(t)

		q, err := query.New("my-query", "select * from transactions")
		gm.Expect(err).To(gm.BeNil())

		evaluator := NewEvaluator(q, schema, make([]*primitives.Policy, 0)...)
		_, err = evaluator.Evaluate()
		gm.Expect(err.Error()).To(gm.Equal("access_denied: No policies match the provided query"))
	})

	t.Run("Fails if there are policies attached but none match", func(t *testing.T) {
		gm.RegisterTestingT(t)

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.othertable",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{primitives.Star},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "select * from transactions")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		_, err = e.Evaluate()
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Passes when * is allowed", func(t *testing.T) {
		gm.RegisterTestingT(t)

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{primitives.Star},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "select * from transactions")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		_, err = e.Evaluate()
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Passes when you are specific about fields", func(t *testing.T) {
		gm.RegisterTestingT(t)

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{"card_number", "vendor"},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "select card_number, vendor from transactions")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		_, err = e.Evaluate()
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Fails when you are specific about fields and request something you cannot see", func(t *testing.T) {
		gm.RegisterTestingT(t)

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{"card_number", "vendor"},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "select card_number, vendor, ssn from transactions")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		_, err = e.Evaluate()
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Fails when you ask for a blacklisted field", func(t *testing.T) {
		gm.RegisterTestingT(t)

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{primitives.Star},
				},

				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Deny,
					Fields: []primitives.Field{"card_number"},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "select card_number from transactions")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		_, err = e.Evaluate()
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("rewrites the query if you use a star command", func(t *testing.T) {
		gm.RegisterTestingT(t)

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{"card_number", "processor"},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "SELECT * FROM transactions")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		q, err = e.Evaluate()
		gm.Expect(err).To(gm.BeNil())

		raw, _ := q.Raw()
		gm.Expect(raw).To(gm.Equal("SELECT card_number, processor FROM transactions"))
	})

	t.Run("cannot run a where if it is not allowed", func(t *testing.T) {
		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Deny,
					Fields: []primitives.Field{primitives.Star},
				},

				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{"card_number"},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "SELECT card_number FROM transactions where processor = 'visa'")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		_, err = e.Evaluate()
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("can run a where if it is allowed", func(t *testing.T) {
		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Deny,
					Fields: []primitives.Field{primitives.Star},
				},

				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{"card_number", "processor"},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "SELECT card_number FROM transactions where processor = 'visa'")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		_, err = e.Evaluate()
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can filter rows with where rules", func(t *testing.T) {
		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   "my-policy",
			Rules: []*primitives.Rule{
				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Allow,
					Fields: []primitives.Field{primitives.Star},
				},

				{
					Target: "records:mycollection.transactions",
					Action: primitives.Read,
					Effect: primitives.Deny,
					Where: []primitives.Where{
						{"processor": "visa"},
					},
				},
			},
		}

		p, err := primitives.NewPolicy("my-policy", spec)
		gm.Expect(err).To(gm.BeNil())

		q, err := query.New("my-query", "SELECT card_number FROM transactions")
		gm.Expect(err).To(gm.BeNil())

		e := NewEvaluator(q, schema, p)

		q, err = e.Evaluate()
		gm.Expect(err).To(gm.BeNil())

		raw, params := q.Raw()
		gm.Expect(raw).To(gm.Equal("SELECT card_number FROM transactions WHERE processor != ?"))
		gm.Expect(params[0]).To(gm.Equal("visa"))
	})
}
