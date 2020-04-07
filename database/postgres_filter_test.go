package database

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

type testStringer struct{}

func (t *testStringer) String() string {
	return "a"
}

func TestPostgresFilter(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := map[string]struct {
		filter  Filter
		success bool
		cause   errors.Cause
		output  string
		values  []interface{}
	}{
		"simpe value comparison": {
			filter:  Filter{Where: Where{"a": "b"}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' = $1",
			values:  []interface{}{"b"},
		},
		"using a Stringer as a value": {
			filter:  Filter{Where: Where{"a": &testStringer{}}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' = $1",
			values:  []interface{}{"a"},
		},
		"two value comparison": {
			filter:  Filter{Where: Where{"a": "b", "c": "d"}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' = $1 AND data::jsonb#>>'{c}' = $2",
			values:  []interface{}{"b", "d"},
		},
		"an In operator": {
			filter:  Filter{Where: Where{"a": In{"a", "b"}}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' IN ($1, $2)",
			values:  []interface{}{"a", "b"},
		},
		"simple comparison AND an IN operator": {
			filter:  Filter{Where: Where{"a": "b", "b": In{"d", "e"}}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' = $1 AND data::jsonb#>>'{b}' IN ($2, $3)",
			values:  []interface{}{"b", "d", "e"},
		},
		"simple comparison with a page no offset": {
			filter:  Filter{Where: Where{"a": "b"}, Page: &Page{Limit: 10}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' = $1 LIMIT 10",
			values:  []interface{}{"b"},
		},
		"simpe comparison with a page with offset": {
			filter:  Filter{Where: Where{"a": "c"}, Page: &Page{Limit: 10, Offset: 2}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' = $1 LIMIT 10 OFFSET 2",
			values:  []interface{}{"c"},
		},
		"simple comparison with an order by": {
			filter:  Filter{Where: Where{"a": "d"}, Order: &Order{Desc, "b"}},
			success: true,
			output:  "WHERE data::jsonb#>>'{a}' = $1 ORDER BY (data::jsonb#>'{b}') DESC",
			values:  []interface{}{"d"},
		},
		"simple comparison with an order by asc": {
			filter:  Filter{Where: Where{"d": "e"}, Order: &Order{Asc, "d"}},
			success: true,
			output:  "WHERE data::jsonb#>>'{d}' = $1 ORDER BY (data::jsonb#>'{d}') ASC",
			values:  []interface{}{"e"},
		},
		"all the components": {
			filter:  Filter{Where: Where{"d": "b"}, Page: &Page{10, 2}, Order: &Order{Asc, "f"}},
			success: true,
			output:  "WHERE data::jsonb#>>'{d}' = $1 ORDER BY (data::jsonb#>'{f}') ASC LIMIT 10 OFFSET 2",
			values:  []interface{}{"b"},
		},
		"errors on an empty In": {
			filter:  Filter{Where: Where{"id": In{}}},
			success: false,
			cause:   NotFoundCause,
		},
	}

	for str, tc := range tests {
		t.Run(str, func(t *testing.T) {
			sql, values, err := buildFilter(tc.filter)
			if !tc.success {
				gm.Expect(errors.FromCause(err, NotFoundCause)).To(gm.BeTrue())
				return
			}

			gm.Expect(sql).To(gm.Equal(tc.output))
			gm.Expect(values).To(gm.Equal(tc.values))
		})
	}
}
