package primitives

// User represents a user of the system
type User struct {
	*Primitive
	Name string `json:"name"`
}

// NewUser returns a new User struct
func NewUser(name string) (*User, error) {
	p, err := newPrimitive(UserType)
	if err != nil {
		return nil, err
	}

	return &User{
		Primitive: p,
		Name:      name, // TODO: Figure out what to do about validation
	}, nil
}
