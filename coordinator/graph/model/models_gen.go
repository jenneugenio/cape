// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

type AttemptRecoveryRequest struct {
	NewPassword primitives.Password `json:"new_password"`
	Secret      primitives.Password `json:"secret"`
	ID          database.ID         `json:"id"`
}

type CreateProjectRequest struct {
	Name        models.ProjectDisplayName `json:"name"`
	Label       *models.Label             `json:"label"`
	Description models.ProjectDescription `json:"Description"`
}

type CreateRecoveryRequest struct {
	Email models.Email `json:"email"`
}

type CreateTokenRequest struct {
	UserID string `json:"user_id"`
}

type CreateTokenResponse struct {
	Secret primitives.Password `json:"secret"`
	Token  *models.Token       `json:"token"`
}

type CreateUserRequest struct {
	Name  models.Name  `json:"name"`
	Email models.Email `json:"email"`
}

type CreateUserResponse struct {
	Password primitives.Password `json:"password"`
	User     *models.User        `json:"user"`
}

type DeleteRecoveriesRequest struct {
	Ids []database.ID `json:"ids"`
}

type ProjectSpecFile struct {
	Transformations []*models.NamedTransformation `json:"transformations"`
	Rules           []*models.Rule                `json:"rules"`
}

type UpdateProjectRequest struct {
	Name        *models.ProjectDisplayName `json:"name"`
	Description *models.ProjectDescription `json:"description"`
}
