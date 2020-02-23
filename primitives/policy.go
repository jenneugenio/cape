package primitives

// Policy is a single defined policy
// TODO -- write this
type Policy struct {
	*Primitive
}

// NewPolicy returns a mutable policy struct
func NewPolicy() (*Policy, error) {
	p, err := newPrimitive(PolicyType)
	if err != nil {
		return nil, err
	}

	return &Policy{
		Primitive: p,
	}, nil
}
