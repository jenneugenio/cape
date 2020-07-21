// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

type ActionInput struct {
	Transform models.Transformation `json:"transform"`
}

type AssignRoleRequest struct {
	RoleID database.ID `json:"role_id"`
	UserID string      `json:"user_id"`
}

type Assignment struct {
	ID        database.ID      `json:"id"`
	Role      *primitives.Role `json:"role"`
	User      *models.User     `json:"user"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type AttachPolicyRequest struct {
	PolicyID string      `json:"policy_id"`
	RoleID   database.ID `json:"role_id"`
}

type Attachment struct {
	ID        database.ID      `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Role      *primitives.Role `json:"role"`
	Policy    *models.Policy   `json:"policy"`
}

type AttemptRecoveryRequest struct {
	NewPassword primitives.Password `json:"new_password"`
	Secret      primitives.Password `json:"secret"`
	ID          database.ID         `json:"id"`
}

type CreatePolicyRequest struct {
	Label primitives.Label `json:"label"`
	Rules []*RuleInput     `json:"rules"`
}

type CreateProjectRequest struct {
	Name        models.ProjectDisplayName `json:"name"`
	Label       *models.Label             `json:"label"`
	Description models.ProjectDescription `json:"Description"`
}

type CreateRecoveryRequest struct {
	Email models.Email `json:"email"`
}

type CreateRoleRequest struct {
	Label   primitives.Label `json:"label"`
	UserIds []string         `json:"user_ids"`
}

type CreateTokenRequest struct {
	UserID string `json:"user_id"`
}

type CreateTokenResponse struct {
	Secret primitives.Password `json:"secret"`
	Token  *primitives.Token   `json:"token"`
}

type CreateUserRequest struct {
	Name  models.Name  `json:"name"`
	Email models.Email `json:"email"`
}

type CreateUserResponse struct {
	Password primitives.Password `json:"password"`
	User     *models.User        `json:"user"`
}

type DeletePolicyRequest struct {
	Label string `json:"label"`
}

type DeleteRecoveriesRequest struct {
	Ids []database.ID `json:"ids"`
}

type DeleteRoleRequest struct {
	ID database.ID `json:"id"`
}

type DetachPolicyRequest struct {
	PolicyID string      `json:"policy_id"`
	RoleID   database.ID `json:"role_id"`
}

type MatchInput struct {
	Name string `json:"name"`
}

type RuleInput struct {
	Match   *MatchInput    `json:"match"`
	Actions []*ActionInput `json:"actions"`
}

type UpdateProjectRequest struct {
	Name        *models.ProjectDisplayName `json:"name"`
	Description *models.ProjectDescription `json:"description"`
}
