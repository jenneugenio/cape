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

func (q *Query) Raw() string {

	selectClause := ""
	for i, f := range q.safe.Fields {
		selectClause += fmt.Sprintf("%s", f)
		if i < len(q.safe.Fields)-1 {
			selectClause += ", "
		}
	}

	return fmt.Sprintf("SELECT %s FROM %s", selectClause, q.safe.TableName)
}

func (q *Query) Rewrite(p *Policy) (*Query, error) {

	unsafe := q.unsafe
	safe := &qq.Query{
		Type:      qq.Select,
		TableName: q.unsafe.TableName,
	}

	for _, rule := range p.rules {

		// skip if this rule doesn't apply
		if rule.target != unsafe.TableName {
			continue
		}

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
	}

	q.safe = safe
	return q, nil
}

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

type Effect int

const (
	Allow Effect = iota
	Deny
)

type Rule struct {
	target string
	effect Effect
	fields []string
}

func NewRule(target string, effect Effect, fields ...string) (*Rule, error) {
	return &Rule{
		target,
		effect,
		fields,
	}, nil
}

type Policy struct {
	rules []*Rule
}

func NewPolicy(rules ...*Rule) (*Policy, error) {
	return &Policy{rules}, nil
}
