// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
	"github.com/manifoldco/go-base64"
)

type AddSourceRequest struct {
	Label       primitives.Label `json:"label"`
	Credentials primitives.DBURL `json:"credentials"`
	ServiceID   *database.ID     `json:"service_id"`
}

type AssignRoleRequest struct {
	RoleID     database.ID `json:"role_id"`
	IdentityID database.ID `json:"identity_id"`
}

type Assignment struct {
	ID        database.ID         `json:"id"`
	Role      *primitives.Role    `json:"role"`
	Identity  primitives.Identity `json:"identity"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type AttachPolicyRequest struct {
	PolicyID database.ID `json:"policy_id"`
	RoleID   database.ID `json:"role_id"`
}

type Attachment struct {
	ID        database.ID        `json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Role      *primitives.Role   `json:"role"`
	Policy    *primitives.Policy `json:"policy"`
}

type CreatePolicyRequest struct {
	Label primitives.Label      `json:"label"`
	Spec  primitives.PolicySpec `json:"spec"`
}

type CreateProjectRequest struct {
	Name        primitives.DisplayName `json:"name"`
	Label       *primitives.Label      `json:"label"`
	Description primitives.Description `json:"Description"`
}

type CreateProjectSpecRequest struct {
	ID        database.ID              `json:"id"`
	ProjectID database.ID              `json:"project_id"`
	ParentID  *database.ID             `json:"parent_id"`
	SourceIds []database.ID            `json:"source_ids"`
	Policies  []*primitives.PolicySpec `json:"policies"`
}

type CreateRoleRequest struct {
	Label       primitives.Label `json:"label"`
	IdentityIds []database.ID    `json:"identity_ids"`
}

type CreateServiceRequest struct {
	Email    primitives.Email       `json:"email"`
	Type     primitives.ServiceType `json:"type"`
	Endpoint *primitives.URL        `json:"endpoint"`
}

type CreateTokenRequest struct {
	IdentityID database.ID `json:"identity_id"`
}

type CreateTokenResponse struct {
	Secret primitives.Password `json:"secret"`
	Token  *primitives.Token   `json:"token"`
}

type CreateUserRequest struct {
	Name  primitives.Name  `json:"name"`
	Email primitives.Email `json:"email"`
}

type CreateUserResponse struct {
	Password primitives.Password `json:"password"`
	User     *primitives.User    `json:"user"`
}

type DeletePolicyRequest struct {
	ID database.ID `json:"id"`
}

type DeleteRoleRequest struct {
	ID database.ID `json:"id"`
}

type DeleteServiceRequest struct {
	ID database.ID `json:"id"`
}

type DeleteSessionRequest struct {
	Token *base64.Value `json:"token"`
}

type DetachPolicyRequest struct {
	PolicyID database.ID `json:"policy_id"`
	RoleID   database.ID `json:"role_id"`
}

type PolicyInput struct {
	Label primitives.Label `json:"label"`
}

type RemoveSourceRequest struct {
	Label primitives.Label `json:"label"`
}

type ReportSchemaRequest struct {
	SourceID     database.ID `json:"source_id"`
	SourceSchema string      `json:"source_schema"`
}

type SessionRequest struct {
	Email   *primitives.Email   `json:"email"`
	TokenID *database.ID        `json:"token_id"`
	Secret  primitives.Password `json:"secret"`
}

type SetupRequest struct {
	Name     primitives.Name     `json:"name"`
	Email    primitives.Email    `json:"email"`
	Password primitives.Password `json:"password"`
}

type UpdateProjectRequest struct {
	Name        *primitives.DisplayName `json:"name"`
	Label       *primitives.Label       `json:"label"`
	Description *primitives.Description `json:"Description"`
}

type UpdateSourceRequest struct {
	SourceLabel primitives.Label `json:"source_label"`
	ServiceID   *database.ID     `json:"service_id"`
}
