package database

// UnsupportedBackendError is raised when you attempt to create a db backened that we do not support
type UnsupportedBackendError struct {
	tried string
}

func (e *UnsupportedBackendError) Error() string {
	return "Attempted to start a " + e.tried + " backend."
}

func newUnsupportedBackendError(tried string) error {
	return &UnsupportedBackendError{tried}
}
