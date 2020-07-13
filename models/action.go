package models

// RBACAction represents what kinds of actions this policy applies to
type RBACAction string

const (
	Create RBACAction = "create"
	Read   RBACAction = "read"
	Update RBACAction = "update"
	Delete RBACAction = "delete"
)
