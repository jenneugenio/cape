package policy

import (
	"github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/dropoutlabs/cape/query"
)

// Evaluator takes a query, schema, and multiple policies and then evaluates a query
// Either modifying the query if it makes sense to do so, or returning an error
type Evaluator struct {
	q      *query.Query
	s      *proto.Schema
	allows []*primitives.Rule
	denies []*primitives.Rule
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []primitives.Field) []primitives.Field {
	mb := make(map[primitives.Field]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []primitives.Field
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func getFields(rules []*primitives.Rule) []primitives.Field {
	var fields []primitives.Field
	for _, r := range rules {
		for _, f := range r.Fields {
			contained := false
			for _, of := range fields {
				if f == of {
					contained = true
				}
			}

			if !contained {
				fields = append(fields, f)
			}
		}
	}

	return fields
}

// Does any policy all everything?
func (e *Evaluator) ruleAllowsAll() bool {
	for _, r := range e.allows {
		for _, f := range r.Fields {
			if f == primitives.Star {
				return true
			}
		}
	}

	return false
}

func (e *Evaluator) deniedFields() []primitives.Field {
	return getFields(e.denies)
}

func (e *Evaluator) allowedFields() []primitives.Field {
	return getFields(e.allows)
}

func (e *Evaluator) evaluateStar() (*query.Query, error) {
	// if they want everything, we need to find a rule that says they can access everything
	// OR, we need to find out what they can access

	var fields []primitives.Field
	if e.ruleAllowsAll() {
		for _, f := range e.s.Fields {
			fields = append(fields, primitives.Field(f.String()))
		}
	} else {
		fields = e.allowedFields()
	}

	if len(fields) == 0 {
		return nil, errors.New(AccessDeniedCause, "No policies match the provided query")
	}

	// now, remove the fields that our policy denies
	fields = difference(fields, e.deniedFields())
	e.q.SetFields(fields)
	return e.q, nil
}

// Evaluate the provided query against the provided policies
// This will return an error if the query violates policy, or return a version
// of the query that is safe to run
func (e *Evaluator) Evaluate() (*query.Query, error) {
	if e.q.WantStar() {
		return e.evaluateStar()
	}

	if len(e.allows) == 0 {
		return nil, errors.New(AccessDeniedCause, "No policies match the provided query")
	}

	// Now, we must find a rule that allows the provided query to run
	// There are basically two cases here
	// The rule could specify a wildcard (*) yes, in which case the query is allowed to run
	allowed := false

	// The second case is that the rules are more specific, in that case we need to find a rule that lets
	// lets each requested target run
	requestedFields := map[primitives.Field]bool{}
	for _, f := range e.q.Fields() {
		requestedFields[f] = false
	}

	for _, r := range e.allows {
		fields := r.Fields
		for _, f := range fields {
			// special case -- the rule says they can access anything
			if f == primitives.Star {
				allowed = true
			}

			requestedFields[f] = true
		}
	}

	if !allowed {
		allTrue := true
		for _, v := range requestedFields {
			if !v {
				allTrue = false
			}
		}

		if allTrue {
			allowed = true
		}
	}

	if !allowed {
		return nil, errors.New(AccessDeniedCause, "No policies allow the requested action to run")
	}

	denied := false
	for _, r := range e.denies {
		fields := r.Fields
		for _, f := range fields {
			for _, qf := range e.q.Fields() {
				if qf == f {
					denied = true
				}
			}
		}
	}

	if denied {
		return nil, errors.New(AccessDeniedCause, "Policy denies the query from running")
	}

	return e.q, nil
}

// New returns a new Evaluator
func New(q *query.Query, s *proto.Schema, policies ...*primitives.Policy) *Evaluator {
	// find the policies that target the given query
	var allows []*primitives.Rule
	var denies []*primitives.Rule
	for _, p := range policies {
		for _, r := range p.Spec.Rules {
			if r.Target.Entity().String() == q.Entity() && r.Effect == primitives.Allow {
				allows = append(allows, r)
			}

			if r.Target.Entity().String() == q.Entity() && r.Effect == primitives.Deny {
				denies = append(denies, r)
			}
		}
	}

	return &Evaluator{
		allows: allows,
		denies: denies,
		s:      s,
		q:      q,
	}
}
