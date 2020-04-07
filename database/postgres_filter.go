package database

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dropoutlabs/cape/database/types"
)

// buildInsert returns the values statement for a multi-value query
func buildInsert(entities []Entity, t types.Type) (string, []interface{}) {
	rows := []string{}
	values := []interface{}{}
	for i, e := range entities {
		rows = append(rows, fmt.Sprintf("($%d)", i+1))
		values = append(values, e)

		if e.GetType() != t {
			panic("not all entities match type: " + t.String())
		}
	}

	return strings.Join(rows, ", "), values
}

// buildFilter returns the clause and parameters for a postgres query
func buildFilter(f Filter) (string, []interface{}, error) {
	fields := []string{}
	values := []interface{}{}
	count := 1 // We need to number the parameters

	// We need to sort the keys in the map in order to build the same SQL
	// statement everytime. This is generally unnecessary but is required for
	// testing scenarios.
	keys := []string{}
	for k := range f.Where {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// Build up the list of values & statemen
	for _, key := range keys {
		value := f.Where[key]
		path := buildDataPath(key, true)
		switch item := value.(type) {
		case In:
			if item.Empty() {
				return "", nil, ErrEmptyIn
			}

			// Compute the results of the IN Clause
			results := item.Uniqify().Values()
			values = append(values, results...)
			params := createParams(count, len(results))

			field := fmt.Sprintf("%s IN (%s)", path, strings.Join(params, ", "))
			fields = append(fields, field)
			count += len(results)
		case fmt.Stringer:
			fields = append(fields, fmt.Sprintf("%s = $%d", path, count))
			values = append(values, item.String())
			count++
		default:
			fields = append(fields, fmt.Sprintf("%s = $%d", path, count))
			values = append(values, item)
			count++
		}
	}

	out := ""
	if len(fields) > 0 {
		out = fmt.Sprintf("WHERE %s", strings.Join(fields, " AND "))
	}

	if f.Order != nil {
		field := buildDataPath(f.Order.Field, false)
		dir := "DESC"
		if f.Order.Dir == Asc {
			dir = "ASC"
		}

		out = fmt.Sprintf("%s ORDER BY (%s) %s", out, field, dir)
	}

	if f.Page != nil {
		lim := ""
		off := ""
		if f.Page.Limit != 0 {
			lim = fmt.Sprintf(" LIMIT %d", f.Page.Limit)
		}

		if f.Page.Offset != 0 {
			off = fmt.Sprintf(" OFFSET %d", f.Page.Offset)
		}

		out = fmt.Sprintf("%s%s%s", out, lim, off)
	}

	return out, values, nil
}

// buildDataPath returns the path to the given field in postgres
func buildDataPath(path string, asText bool) string {
	selector := ">>"
	if !asText {
		selector = ">"
	}

	path = strings.Replace(path, ".", ",", -1)
	return fmt.Sprintf("data::jsonb#%s'{%s}'", selector, path)
}

func createParams(start, n int) []string {
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = fmt.Sprintf("$%d", start+i)
	}

	return out
}
