package partyerrors

// List of all possible causes
var (
	UnsupportedErrorCause = NewCause(InternalServerErrorCategory, "unsupported_error")
	UnknownCause          = NewCause(InternalServerErrorCategory, "unknown_cause")
	InvalidArgumentCause  = NewCause(BadRequestCategory, "invalid_argument")
	InvalidStateCause     = NewCause(BadRequestCategory, "invalid_state")
	TimeoutCause          = NewCause(RequestTimeoutCategory, "timeout")
	NotImplementedCause   = NewCause(NotImplementedCategory, "not_implemented")
)

// Cause is the cause of an error
type Cause struct {
	Name     string
	Category Category
}

// NewCause creates a new cause of an error
func NewCause(c Category, name string) Cause {
	return Cause{
		Category: c,
		Name:     name,
	}
}

// CausedBy check if an error has a given cause
func CausedBy(err error, c Cause) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.Cause == c
}
