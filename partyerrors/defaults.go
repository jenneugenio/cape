package partyerrors

var (
	ErrNotImplemented = New(NotImplementedCause, "Not Implemented")
)
