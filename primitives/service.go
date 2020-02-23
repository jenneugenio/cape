package primitives

// Service represents a service connecting to the system (e.g. a machine
// running a pipeline).
type Service struct {
	*Primitive
}

// NewService returns a mutable service struct
func NewService() (*Service, error) {
	p, err := newPrimitive(ServiceType)
	if err != nil {
		return nil, err
	}

	return &Service{
		Primitive: p,
	}, nil
}
