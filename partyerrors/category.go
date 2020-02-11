package partyerrors

import (
	"fmt"
	"net/http"
)

// Category is a category of cause
type Category string

// String returns the stringified name of the category
func (t Category) String() string {
	return string(t)
}

// List of all error categories
var (
	BadRequestCategory          Category = "bad_request"
	UnauthorizedCategory        Category = "unauthorized"
	NotFoundCategory            Category = "not_found"
	InternalServerErrorCategory Category = "internal_error"
	NotImplementedCategory      Category = "not_implemented"
	MethodNotAllowedCategory    Category = "method_not_allowed"
	ConflictCategory            Category = "conflict"
	ForbiddenCategory           Category = "forbidden"
	RequestTimeoutCategory      Category = "request_timeout"
)

func init() {
	for k, v := range statusCodeCategoryMap {
		if _, ok := inverseMap[v]; ok {
			panic(fmt.Sprintf("The status code %d is already bound to the category named: %s", v, k))
		}

		inverseMap[v] = k
	}
}

var statusCodeCategoryMap = map[Category]int{
	BadRequestCategory:          http.StatusBadRequest,
	UnauthorizedCategory:        http.StatusUnauthorized,
	NotFoundCategory:            http.StatusNotFound,
	InternalServerErrorCategory: http.StatusInternalServerError,
	NotImplementedCategory:      http.StatusNotImplemented,
	MethodNotAllowedCategory:    http.StatusMethodNotAllowed,
	ConflictCategory:            http.StatusConflict,
	ForbiddenCategory:           http.StatusForbidden,
	RequestTimeoutCategory:      http.StatusRequestTimeout,
}

var inverseMap = map[int]Category{}

func statusCodeForCategory(c Category) (int, bool) {
	code, ok := statusCodeCategoryMap[c]
	return code, ok
}

// CategoryForStatusCode convert an http status code to a category
func CategoryForStatusCode(code int) (Category, bool) {
	t, ok := inverseMap[code]
	return t, ok
}
