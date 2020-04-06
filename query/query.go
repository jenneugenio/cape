package query

import (
	"fmt"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/marianogappa/sqlparser"
	qq "github.com/marianogappa/sqlparser/query"
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

// Entity returns what this query is targeting (e.g. a postgres table)
func (q *Query) Entity() string {
	return q.q.TableName
}

func (q *Query) Fields() []primitives.Field {
	fields := make([]primitives.Field, len(q.q.Fields))
	for i, f := range q.q.Fields {
		fields[i] = primitives.Field(f)
	}

	return fields
}

func (q *Query) SetFields(fields []primitives.Field) {
	fStr := make([]string, len(fields))
	for i, f := range fields {
		fStr[i] = f.String()
	}

	q.q.Fields = fStr
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

func (q *Query) WantStar() bool {
	for _, f := range q.q.Fields {
		if f == "*" {
			return true
		}
	}

	return false
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
