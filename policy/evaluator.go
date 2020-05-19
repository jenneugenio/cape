package policy

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/capeprivacy/cape/query"
)

// Evaluator takes a query, schema, and multiple policies and then evaluates a query
// Either modifying the query if it makes sense to do so, or returning an error
type Evaluator struct {
	q               *query.Query
	s               *proto.Schema
	allowFieldRules []*primitives.Rule
	denyFieldRules  []*primitives.Rule
	allowWhereRules []*primitives.Rule
	denyWhereRules  []*primitives.Rule

	transforms []*primitives.Transformation
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
	for _, r := range e.allowFieldRules {
		for _, f := range r.Fields {
			if f == primitives.Star {
				return true
			}
		}
	}

	return false
}

func (e *Evaluator) deniedFields() []primitives.Field {
	return getFields(e.denyFieldRules)
}

func (e *Evaluator) allowedFields() []primitives.Field {
	return getFields(e.allowFieldRules)
}

func (e *Evaluator) attachConditions() error {
	if len(e.allowWhereRules) > 0 {
		var where []primitives.Where
		for _, r := range e.allowWhereRules {
			where = append(where, r.Where...)
		}

		e.q.SetConditions(where, primitives.Eq)
	}

	if len(e.denyWhereRules) > 0 {
		var where []primitives.Where
		for _, r := range e.denyWhereRules {
			where = append(where, r.Where...)
		}

		e.q.SetConditions(where, primitives.Neq)
	}

	return nil
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
		return nil, errors.New(auth.AuthorizationFailure, "No policies match the provided query")
	}

	if e.checkDeniedWheres() {
		return nil, errors.New(auth.AuthorizationFailure, "Policy where rules denies the query from running")
	}

	// now, remove the fields that our policy denyFieldRules
	fields = difference(fields, e.deniedFields())
	e.q.SetFields(fields)
	err := e.attachConditions()
	if err != nil {
		return nil, err
	}

	return e.q, nil
}

// Evaluate the provided query against the provided policies
// This will return an error if the query violates policy, or return a version
// of the query that is safe to run
func (e *Evaluator) Evaluate() (*query.Query, error) {
	if e.q.WantStar() {
		return e.evaluateStar()
	}

	if len(e.allowFieldRules) == 0 {
		return nil, errors.New(auth.AuthorizationFailure, "No policies match the provided query")
	}

	// Now, we must find a rule that allowFieldRules the provided query to run
	// There are basically two cases here
	// The rule could specify a wildcard (*) yes, in which case the query is allowed to run
	allowed := false

	// The second case is that the rules are more specific, in that case we need to find a rule that lets
	// lets each requested target run
	requestedFields := map[primitives.Field]bool{}
	for _, f := range e.q.Fields() {
		requestedFields[f] = false
	}

	// also check which fields they are using in a conditional
	for f := range e.q.Conditions() {
		requestedFields[f] = false
	}

	for _, r := range e.allowFieldRules {
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
		return nil, errors.New(auth.AuthorizationFailure, "No policies allow the requested action to run")
	}

	denied := false
	for _, r := range e.denyFieldRules {
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
		return nil, errors.New(auth.AuthorizationFailure, "Policy denyFieldRules the query from running")
	}

	if e.checkDeniedWheres() {
		return nil, errors.New(auth.AuthorizationFailure, "Policy where rules denies the query from running")
	}

	err := e.attachConditions()
	if err != nil {
		return nil, err
	}

	return e.q, nil
}

// Transforms returns all the transforms contained in the policies
func (e *Evaluator) Transforms() []*primitives.Transformation {
	return e.transforms
}

func (e *Evaluator) checkDeniedWheres() bool {
	// This checks to see if there is a condition in the submitted query
	// that has been disallowed. If one is found it denies the whole query.
	// For example, if a policy does not allow access to a "processor" field
	// where it equals "Visa" and a user submits a query like:
	// "SELECT * FROM transactions WHERE processor = 'Visa'" then the
	// query will be rejected.
	denied := false
	for _, r := range e.denyWhereRules {
		wheres := r.Where
		for _, w := range wheres {
			for f, qw := range e.q.Conditions() {
				if qw == w[f.String()] {
					denied = true
				}
			}
		}
	}

	return denied
}

// NewEvaluator returns a new Evaluator
func NewEvaluator(q *query.Query, s *proto.Schema, policies ...*primitives.Policy) *Evaluator {
	// find the policies that target the given query
	var allowFieldRules []*primitives.Rule
	var denyFieldRules []*primitives.Rule
	var allowWhereRules []*primitives.Rule
	var denyWhereRules []*primitives.Rule
	var transforms []*primitives.Transformation

	for _, p := range policies {
		for _, r := range p.Spec.Rules {
			transforms = append(transforms, r.Transformations...)

			if r.Type() == primitives.FieldRule && r.Target.Entity().String() == q.Entity() && r.Effect == primitives.Allow {
				allowFieldRules = append(allowFieldRules, r)
			}

			if r.Type() == primitives.FieldRule && r.Target.Entity().String() == q.Entity() && r.Effect == primitives.Deny {
				denyFieldRules = append(denyFieldRules, r)
			}

			if r.Type() == primitives.WhereRule && r.Target.Entity().String() == q.Entity() && r.Effect == primitives.Allow {
				allowWhereRules = append(allowWhereRules, r)
			}

			if r.Type() == primitives.WhereRule && r.Target.Entity().String() == q.Entity() && r.Effect == primitives.Deny {
				denyWhereRules = append(denyWhereRules, r)
			}
		}
	}

	return &Evaluator{
		allowFieldRules: allowFieldRules,
		denyFieldRules:  denyFieldRules,
		allowWhereRules: allowWhereRules,
		denyWhereRules:  denyWhereRules,
		transforms:      transforms,
		s:               s,
		q:               q,
	}
}
