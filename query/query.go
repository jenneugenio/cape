package query

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/marianogappa/sqlparser"
	qq "github.com/marianogappa/sqlparser/query"
)

// Query represents a data query that should run against a source's collection
type Query struct {
	unsafe *qq.Query
	source primitives.Label
}

// Source returns the data source that this query should run against (e.g. which database)
func (q *Query) Source() primitives.Label {
	return q.source
}

// Collection returns which data collection this query should run against (e.g. which table)
func (q *Query) Collection() string {
	return q.unsafe.TableName
}

// Validate the query
func (q *Query) Validate() error {
	if q.unsafe.Type != qq.Select {
		return errors.New(InvalidQueryCause, "Only select statements are supported")
	}

	return nil
}

// Target returns what this query is targeting (e.g. a postgres table)
func (q *Query) Target() string {
	return q.unsafe.TableName
}

// New creates a new query object
func New(source primitives.Label, query string) (*Query, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	q := &Query{
		unsafe: &stmt,
		source: source,
	}

	err = q.Validate()
	if err != nil {
		return nil, err
	}

	return q, nil
}
