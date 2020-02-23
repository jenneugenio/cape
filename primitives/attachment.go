package primitives

// Attachment represents a policy being applied/attached to a role
type Attachment struct {
	*Primitive
}

// NewAttachment returns a new attachment
func NewAttachment() (*Attachment, error) {
	p, err := newPrimitive(AttachmentType)
	if err != nil {
		return nil, err
	}

	// TODO: Pass in the values of the attachment!
	//
	// An attachment is considered an immutable type in our object system (as
	// defined by the type)
	a := &Attachment{
		Primitive: p,
	}

	ID, err := DeriveID(a)
	if err != nil {
		return nil, err
	}

	a.ID = ID
	return a, nil
}
