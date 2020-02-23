package primitives

// Role in a role in the system (e.g. Admin, user, etc)
type Role struct {
	*Primitive
}

// NewRole returns a mutable role struct
func NewRole() (*Role, error) {
	p, err := newPrimitive(RoleType)
	if err != nil {
		return nil, err
	}

	return &Role{
		Primitive: p,
	}, nil
}
