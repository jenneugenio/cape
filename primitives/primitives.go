// Package primitives contains all the primitive data types used by the coordinator
// and the connector.
package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database/types"
)

var (
	// UserType represents a user (ie a person)
	UserType types.Type = 0x000

	// ServicePrimitiveType is a service (e.g. another machine)
	ServicePrimitiveType types.Type = 0x001

	// TokenPrimitiveType is a token representing an authenticated user
	TokenPrimitiveType types.Type = 0x002

	// RoleType represents a role, which is attached to users/services to determine what they can or cannot do
	RoleType types.Type = 0x003

	// PolicyType represents a policy, which dictates how roles must be use
	PolicyType types.Type = 0x004

	// AttachmentType represents a policy being attached to a role
	AttachmentType types.Type = 0x005

	// AssignmentType represents a role being applied to an identity
	AssignmentType types.Type = 0x006

	// SessionType represents a session between a user/services & the system
	SessionType types.Type = 0x007

	// SourcePrimitiveType represents an external database/dataset
	SourcePrimitiveType types.Type = 0x008

	// ConfigType represents the config object for cape
	ConfigType types.Type = 0x009
)

func init() {
	types.Register(UserType, "users", true)
	types.Register(ServicePrimitiveType, "services", true)
	types.Register(TokenPrimitiveType, "tokens", false)
	types.Register(RoleType, "roles", true)
	types.Register(PolicyType, "policies", true)
	types.Register(AttachmentType, "attachments", false)
	types.Register(AssignmentType, "assignments", false)
	types.Register(SessionType, "sessions", false)
	types.Register(SourcePrimitiveType, "sources", true)
	types.Register(ConfigType, "config", true)
}
