package primitives

import (
	"github.com/dropoutlabs/cape/database/types"
)

var (
	// The following types represent the core primitives in the system
	UserType       types.Type = 0x000
	ServiceType    types.Type = 0x001
	TokenType      types.Type = 0x002
	RoleType       types.Type = 0x003
	PolicyType     types.Type = 0x004
	AttachmentType types.Type = 0x005
	AssignmentType types.Type = 0x006
)

func init() {
	types.Register(UserType, "users", false)
	types.Register(ServiceType, "services", false)
	types.Register(TokenType, "tokens", true)
	types.Register(RoleType, "roles", false)
	types.Register(PolicyType, "policies", false)
	types.Register(AttachmentType, "attachments", true)
	types.Register(AssignmentType, "assignments", true)
}
