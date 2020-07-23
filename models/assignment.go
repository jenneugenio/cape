package models

import "time"

// Assignment represents a policy being applied/attached to a role
type Assignment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	ProjectID string    `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Assignment) Validate() error {
	//if err := a.Primitive.Validate(); err != nil {
	//	return errors.Wrap(InvalidAssignmentCause, err)
	//}
	//
	//if a.UserID == "" {
	//	return errors.New(InvalidAssignmentCause, "Invalid Identity ID provided")
	//}
	//
	//if err := a.RoleID.Validate(); err != nil {
	//	return errors.New(InvalidAssignmentCause, "Assignment role id must be valid")
	//}
	//
	//typ, err := a.RoleID.Type()
	//if err != nil {
	//	return errors.New(InvalidAssignmentCause, "Invalid Role ID provider")
	//}
	//
	//if typ != RoleType {
	//	return errors.New(InvalidAssignmentCause, "Invalid Role ID provider")
	//}

	return nil
}

// NewAssignment returns a new Assignment
func NewAssignment(userID string, roleID string) (*Assignment, error) {
	// An Assignment is considered an immutable type in our object system (as
	// defined by the type)
	a := &Assignment{
		UserID: userID,
		RoleID: roleID,
	}

	return a, a.Validate()
}

func (a *Assignment) GetEncryptable() bool {
	return false
}
