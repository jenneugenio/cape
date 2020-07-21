package models

import (
	"fmt"
)

type ProjectStatus string

const (
	ProjectPending  ProjectStatus = "Pending"
	ProjectActive   ProjectStatus = "Active"
	ProjectArchived ProjectStatus = "Archived"
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
}

func NewProject(name ProjectDisplayName, label Label, description ProjectDescription) Project {
	return Project{
		ID:          NewID(),
		Name:        name,
		Label:       label,
		Description: description,
		Status:      ProjectPending,
	}
}
