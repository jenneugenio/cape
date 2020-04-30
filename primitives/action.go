package primitives

// Action represents what kinds of actions this policy applies to
type Action string

const (
	Create Action = "create"
	Read   Action = "read"
	Update Action = "update"
	Delete Action = "delete"
)
