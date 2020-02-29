package database

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestPostgresFilter(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := map[string]struct {
		f Filter
		o string
		v []interface{}
	}{
		"simpe value comparison": {
			Filter{Where: Where{"a": "b"}},
			"WHERE data::jsonb#>>'{a}' = $1",
			[]interface{}{"b"},
		},
		"two value comparison": {
			Filter{Where: Where{"a": "b", "c": "d"}},
			"WHERE data::jsonb#>>'{a}' = $1 AND data::jsonb#>>'{c}' = $2",
			[]interface{}{"b", "d"},
		},
		"an In operator": {
			Filter{Where: Where{"a": In{"a", "b"}}},
			"WHERE data::jsonb#>>'{a}' IN ($1, $2)",
			[]interface{}{"a", "b"},
		},
		"simple comparison AND an IN operator": {
			Filter{Where: Where{"a": "b", "b": In{"d", "e"}}},
			"WHERE data::jsonb#>>'{a}' = $1 AND data::jsonb#>>'{b}' IN ($2, $3)",
			[]interface{}{"b", "d", "e"},
		},
		"simple comparison with a page no offset": {
			Filter{Where: Where{"a": "b"}, Page: &Page{Limit: 10}},
			"WHERE data::jsonb#>>'{a}' = $1 LIMIT 10",
			[]interface{}{"b"},
		},
		"simpe comparison with a page with offset": {
			Filter{Where: Where{"a": "c"}, Page: &Page{Limit: 10, Offset: 2}},
			"WHERE data::jsonb#>>'{a}' = $1 LIMIT 10 OFFSET 2",
			[]interface{}{"c"},
		},
		"simple comparison with an order by": {
			Filter{Where: Where{"a": "d"}, Order: &Order{Desc, "b"}},
			"WHERE data::jsonb#>>'{a}' = $1 ORDER BY (data::jsonb#>'{b}') DESC",
			[]interface{}{"d"},
		},
		"simple comparison with an order by asc": {
			Filter{Where: Where{"d": "e"}, Order: &Order{Asc, "d"}},
			"WHERE data::jsonb#>>'{d}' = $1 ORDER BY (data::jsonb#>'{d}') ASC",
			[]interface{}{"e"},
		},
		"all the components": {
			Filter{Where: Where{"d": "b"}, Page: &Page{10, 2}, Order: &Order{Asc, "f"}},
			"WHERE data::jsonb#>>'{d}' = $1 ORDER BY (data::jsonb#>'{f}') ASC LIMIT 10 OFFSET 2",
			[]interface{}{"b"},
		},
	}

	for str, tc := range tests {
		t.Run(str, func(t *testing.T) {
			sql, values := buildFilter(tc.f)
			gm.Expect(sql).To(gm.Equal(tc.o))
			gm.Expect(values).To(gm.Equal(tc.v))
		})
	}
}
