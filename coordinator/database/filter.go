package database

import (
	"fmt"
	"reflect"
)

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

type Select struct {
	LookingFor string
}

// Filter represents a filter thats applied to a Query
type Filter struct {
	Where Where
	Order *Order
	Page  *Page
}

// In is an operator that requires a field to match
//
// XXX: If we add `OR` support in future then we need to revisit the logic
// we've added around handling an empty `in` operator.
type In []interface{}

// PickerFunc represents a function that will pick a value off an entity
// to be used for performing an In clause as a part of a Filter.
type PickerFunc func(interface{}) interface{}

// InFromEntities returns an In for the given list of entities. If no
// PickerFunc is provided then the ID is used.
func InFromEntities(in interface{}, f PickerFunc) In {
	// We don't know anything about the underlying type because of how slices
	// work inside go. To get around this, we need to use reflect to figure out
	// the underlying type and then adjust accordingly.
	//
	// This is necessary due to the inability to assign a value to a position
	// in a slice of interfaces. Read more here:
	// https://github.com/golang/go/wiki/InterfaceSlice
	inValue := reflect.ValueOf(in)
	if inValue.Kind() != reflect.Slice {
		panic("Expected a slice")
	}

	// We need to get a concrete type so we can check that each item in the
	// slice satisfies the Entity interface.
	entityType := reflect.TypeOf((*Entity)(nil)).Elem()

	len := inValue.Len()
	out := make(In, len)
	for i := 0; i < len; i++ {
		v := inValue.Index(i)
		if !v.Type().Implements(entityType) && v.Kind() == reflect.Ptr {
			v = inValue.Index(i).Elem()
		}

		value := v.Interface().(Entity)
		out[i] = f(value)
	}

	return out
}

// Uniquify returns a copy of In with all duplicate values removed
func (in In) Uniqify() In {
	result := In{}
	found := map[string]bool{}
	for _, value := range in {
		out := ""
		switch v := value.(type) {
		case fmt.Stringer:
			out = v.String()
		case string:
			out = v
		default:
			panic(fmt.Sprintf("In type must be string or Stringer, got %T", v))
		}

		if ok := found[out]; !ok {
			found[out] = true
			result = append(result, value)
		}
	}

	return result
}

// Values returns a slice of strings for all values contained in the In
func (in In) Values() []interface{} {
	result := []interface{}{}
	for _, value := range in {
		switch v := value.(type) {
		case fmt.Stringer:
			result = append(result, v.String())
		case string:
			result = append(result, v)
		default:
			panic(fmt.Sprintf("In type must be string or Stringer, got %T", v))
		}
	}

	return result
}

// Empty returns whether or not the given In contains any actual values
func (in In) Empty() bool {
	return len(in) == 0
}

// NewFilter is a convenience function for creating a Filter
func NewFilter(w Where, o *Order, p *Page) Filter {
	return Filter{Where: w, Order: o, Page: p}
}

// NewEmptyFilter creates an empty filter
func NewEmptyFilter() Filter {
	return Filter{}
}
