package models

import (
	"fmt"
	"time"
)

type ProjectStatus string

const (
	ProjectPending  ProjectStatus = "Pending"
	ProjectActive   ProjectStatus = "Active"
	ProjectArchived ProjectStatus = "Archived"

	Any ProjectStatus = "any"
)

func (p ProjectStatus) String() string {
	return string(p)
}

func (p ProjectStatus) Validate() error {
	switch p {
	case ProjectPending:
		return nil
	case ProjectActive:
		return nil
	case ProjectArchived:
		return nil
	case Any:
		return nil
	}

	return fmt.Errorf("invalid project status: %s", p)
}

type ProjectDescription string

func (p ProjectDescription) String() string {
	return string(p)
}

type ProjectDisplayName string

func (p ProjectDisplayName) String() string {
	return string(p)
}

type Project struct {
	ID            string             `json:"id"`
	Label         Label              `json:"label"`
	Name          ProjectDisplayName `json:"name"`
	Description   ProjectDescription `json:"description"`
	Status        ProjectStatus      `json:"status"`
	CurrentSpecID string
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewProject(name ProjectDisplayName, label Label, description ProjectDescription) Project {
	return Project{
		ID:          NewID(),
		Name:        name,
		Label:       label,
		Description: description,
		Status:      ProjectPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
