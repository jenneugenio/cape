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

type AuthSessionRequest struct {
	Signature base64.Value `json:"signature"`
}

type CreatePolicyRequest struct {
	Label primitives.Label      `json:"label"`
	Spec  primitives.PolicySpec `json:"spec"`
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
	IdentityID database.ID                   `json:"identity_id"`
	PublicKey  base64.Value                  `json:"public_key"`
	Salt       base64.Value                  `json:"salt"`
	Alg        primitives.CredentialsAlgType `json:"alg"`
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

type LoginSessionRequest struct {
	Email   *primitives.Email `json:"email"`
	TokenID *database.ID      `json:"token_id"`
}

type NewUserRequest struct {
	Name      primitives.Name               `json:"name"`
	Email     primitives.Email              `json:"email"`
	PublicKey base64.Value                  `json:"public_key"`
	Salt      base64.Value                  `json:"salt"`
	Alg       primitives.CredentialsAlgType `json:"alg"`
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
