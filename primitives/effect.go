package primitives

// Effect represents what kind of effect this policy has, e.g. allow or deny
type Effect string

const (
	Allow Effect = "allow"
	Deny  Effect = "deny"
)
