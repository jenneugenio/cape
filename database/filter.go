package database

// Direction represents the ordering to be applied to a Filter. The default value
// is descending.
type Direction uint

// A list of the available ordering mechanism
const (
	Desc Direction = iota
	Asc
)

// Where represents a series of conditional clauses applied to a Filter for
// selecting entities based on the value of it's fields.
type Where map[string]interface{}

// Order represents an ordering of fields in a query result
type Order struct {
	Dir   Direction
	Field string
}

// Page represent the filtering properties being applied to only return a
// segment of entities from a query
type Page struct {
	Limit  int
	Offset int
}

// Filter represents a filter thats applied to a Query
type Filter struct {
	Where Where
	Order *Order
	Page  *Page
}

// In is an operator that requires a field to match
type In []interface{}

// NewFilter is a convenience function for creating a Filter
func NewFilter(w Where, o *Order, p *Page) Filter {
	return Filter{Where: w, Order: o, Page: p}
}

// NewEmptyFilter creates an empty filter
func NewEmptyFilter() Filter {
	return Filter{}
}
