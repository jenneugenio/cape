package query

import (
	"fmt"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/marianogappa/sqlparser"
	qq "github.com/marianogappa/sqlparser/query"
	"sort"
)

// Query represents a data query that should run against a source's collection
type Query struct {
	q      *qq.Query
	source primitives.Label
}

// Source returns the data source that this query should run against (e.g. which database)
func (q *Query) Source() primitives.Label {
	return q.source
}

// Collection returns which data collection this query should run against (e.g. which table)
func (q *Query) Collection() string {
	return q.q.TableName
}

// Validate the query
func (q *Query) Validate() error {
	if q.q.Type != qq.Select {
		return errors.New(InvalidQueryCause, "Only select statements are supported")
	}

	return nil
}

// Target returns what this query is targeting (e.g. a postgres table)
func (q *Query) Target() string {
	return q.q.TableName
}

// Raw returns the supplied sql into a policy obeying sql string
func (q *Query) Raw() (string, []interface{}) {
	selectClause := ""
	for i, f := range q.q.Fields {
		selectClause += f
		if i < len(q.q.Fields)-1 {
			selectClause += ", "
		}
	}

	raw := fmt.Sprintf("SELECT %s FROM %s", selectClause, q.q.TableName)
	parameters := make([]interface{}, len(q.q.Conditions))

	if len(q.q.Conditions) > 0 {
		raw += " WHERE "

		for i, c := range q.q.Conditions {
			raw += fmt.Sprintf("%s = ?", c.Operand1)
			parameters[i] = c.Operand2
			if i < len(q.q.Conditions)-1 {
				raw += " AND "
			}
		}
	}

	return raw, parameters
}

func (q *Query) wantStar() bool {
	for _, f := range q.q.Fields {
		if f == "*" {
			return true
		}
	}

	return false
}

// Rewrite the query to something that respects the policy
func (q *Query) Rewrite(p *primitives.Policy) (*Query, error) {
	unsafe := q.q
	safe := &qq.Query{
		Type:      qq.Select,
		TableName: q.q.TableName,
	}

	wantStar := q.wantStar()
	spec := p.Spec

	for _, rule := range spec.Rules {
		// skip if this rule doesn't apply
		if rule.Target.Entity().String() != unsafe.TableName {
			continue
		}

		switch rule.Effect {
		case primitives.Deny:
			// TODO -- This is assuming the rule is to disallow, need to do the opposite if we are allowing
			for _, field := range unsafe.Fields {
				in := false
				for _, denied := range rule.Fields {
					if field == denied.String() {
						in = true
					}
				}

				if !in {
					safe.Fields = append(safe.Fields, field)
				}
			}
		case primitives.Allow:
			for _, allowed := range rule.Fields {
				// If the user has requested *, we will give them everything they can see.  We are short
				// circuiting doing the loop to check if they asked for it, here
				wanted := wantStar
				if !wanted {
					for _, requested := range unsafe.Fields {
						if allowed.String() == requested {
							wanted = true
						}
					}
				}

				if wanted {
					safe.Fields = append(safe.Fields, allowed.String())
				}
			}
		}

		for _, w := range rule.Where {
			// We sort the keys alphabetically on the where clause so that we can ensure their order
			// in the returned statement (e.g. when raw is called).
			// This is useful for testing, ie we can be deterministic w/ regard to the output query
			var keys []string
			for k := range w {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				safe.Conditions = append(safe.Conditions, qq.Condition{
					Operand1:        k,
					Operand1IsField: true,
					Operator:        qq.Eq,
					Operand2:        w[k],
				})
			}
		}
	}

	if len(safe.Fields) == 0 {
		return nil, errors.New(NoPossibleFieldsCause, "Cannot access any requested fields")
	}

	return &Query{
		q:      safe,
		source: q.source,
	}, nil
}

// New creates a new query object
func New(source primitives.Label, query string) (*Query, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	q := &Query{
		q:      &stmt,
		source: source,
	}

	err = q.Validate()
	if err != nil {
		return nil, err
	}

	return q, nil
}
