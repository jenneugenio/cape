package query

import (
	"fmt"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/marianogappa/sqlparser"
	qq "github.com/marianogappa/sqlparser/query"
)

type Query struct {
	unsafe *qq.Query
	safe   *qq.Query
}

// Raw returns the supplied sql into a policy obeying sql string
func (q *Query) Raw() string {
	selectClause := ""
	for i, f := range q.safe.Fields {
		selectClause += f
		if i < len(q.safe.Fields)-1 {
			selectClause += ", "
		}
	}

	raw := fmt.Sprintf("SELECT %s FROM %s", selectClause, q.safe.TableName)
	if len(q.safe.Conditions) > 0 {
		raw += " WHERE "

		for i, c := range q.safe.Conditions {
			raw += fmt.Sprintf("%s = '%s'", c.Operand1, c.Operand2)
			if i < len(q.safe.Conditions)-1 {
				raw += " AND "
			}
		}
	}

	return raw
}

func (q *Query) wantStar() bool {
	for _, f := range q.unsafe.Fields {
		if f == "*" {
			return true
		}
	}

	return false
}

// Rewrite the query to something that respects the policy
func (q *Query) Rewrite(p *Policy) (*Query, error) {
	unsafe := q.unsafe
	safe := &qq.Query{
		Type:      qq.Select,
		TableName: q.unsafe.TableName,
	}

	wantStar := q.wantStar()

	for _, rule := range p.rules {
		// skip if this rule doesn't apply
		if rule.target != unsafe.TableName {
			continue
		}

		switch rule.effect {
		case Deny:
			// TODO -- This is assuming the rule is to disallow, need to do the opposite if we are allowing
			for _, field := range unsafe.Fields {
				in := false
				for _, denied := range rule.fields {
					if field == denied {
						in = true
					}
				}

				if !in {
					safe.Fields = append(safe.Fields, field)
				}
			}
		case Allow:
			for _, allowed := range rule.fields {
				// If the user has requested *, we will give them everything they can see.  We are short
				// circuiting doing the loop to check if they asked for it, here
				wanted := wantStar
				if !wanted {
					for _, requested := range unsafe.Fields {
						if allowed == requested {
							wanted = true
						}
					}
				}

				if wanted {
					safe.Fields = append(safe.Fields, allowed)
				}
			}
		}

		// check conditions on the rule
		for k, v := range rule.condition {
			safe.Conditions = append(safe.Conditions, qq.Condition{
				Operand1:        k,
				Operand1IsField: true,
				Operator:        qq.Eq,
				Operand2:        v,
			})
		}
	}

	if len(safe.Fields) == 0 {
		return nil, errors.New(NoPossibleFieldsCause, "Cannot access any requested fields")
	}

	q.safe = safe
	return q, nil
}

// Parse the incoming query
func Parse(query string) (*Query, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	if stmt.Type != qq.Select {
		return nil, errors.New(InvalidQueryCause, "Only select statements are supported")
	}

	return &Query{
		unsafe: &stmt,
	}, nil
}

// Effect tells whether something can be allowed or denied
type Effect int

const (
	Allow Effect = iota
	Deny
)

// Conditions, e.g. the where in sql
type Condition map[string]string

// Rule represents a single rule within a policy
type Rule struct {
	target    string
	effect    Effect
	fields    []string
	condition Condition
}

// NewRule makes a new Rule
func NewRule(target string, effect Effect, condition Condition, fields ...string) (*Rule, error) {
	return &Rule{
		target,
		effect,
		fields,
		condition,
	}, nil
}

// Policy represents a data policy
type Policy struct {
	rules []*Rule
}

// NewPolicy makes a new policy
func NewPolicy(rules ...*Rule) (*Policy, error) {
	return &Policy{rules}, nil
}
