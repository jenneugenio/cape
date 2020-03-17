package primitives

import (
	"github.com/dropoutlabs/cape/database/types"
)

var (
	// The following types represent the core primitives in the system
	UserType           types.Type = 0x000
	ServiceType        types.Type = 0x001
	TokenPrimitiveType types.Type = 0x002
	RoleType           types.Type = 0x003
	PolicyType         types.Type = 0x004
	AttachmentType     types.Type = 0x005
	AssignmentType     types.Type = 0x006
	SessionType        types.Type = 0x007
)

func init() {
	types.Register(UserType, "users", true)
	types.Register(ServiceType, "services", true)
	types.Register(TokenPrimitiveType, "tokens", false)
	types.Register(RoleType, "roles", true)
	types.Register(PolicyType, "policies", true)
	types.Register(AttachmentType, "attachments", false)
	types.Register(AssignmentType, "assignments", false)
	types.Register(SessionType, "sessions", false)
}
