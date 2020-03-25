// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"net/url"
	"time"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/manifoldco/go-base64"
)

type AddSourceRequest struct {
	Label       primitives.Label `json:"label"`
	Credentials url.URL          `json:"credentials"`
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
	Label primitives.Label `json:"label"`
}

type CreateRoleRequest struct {
	Label       primitives.Label `json:"label"`
	IdentityIds []database.ID    `json:"identity_ids"`
}

type CreateServiceRequest struct {
	Email     primitives.Email              `json:"email"`
	PublicKey base64.Value                  `json:"public_key"`
	Salt      base64.Value                  `json:"salt"`
	Alg       primitives.CredentialsAlgType `json:"alg"`
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
	Token base64.Value `json:"token"`
}

type LoginSessionRequest struct {
	Email primitives.Email `json:"email"`
}

type NewUserRequest struct {
	Name      primitives.Name               `json:"name"`
	Email     primitives.Email              `json:"email"`
	PublicKey base64.Value                  `json:"public_key"`
	Salt      base64.Value                  `json:"salt"`
	Alg       primitives.CredentialsAlgType `json:"alg"`
}
