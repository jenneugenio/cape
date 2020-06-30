package prims

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

const maxDescriptionSize = 5000

type Description string

func (d Description) Validate() error {
	if len(d) > maxDescriptionSize {
		return errors.New(primitives.InvalidProjectDescriptionCause, "%d exceeds max description size of %d", len(d), maxDescriptionSize)
	}

	return nil
}

func (d Description) String() string {
	return string(d)
}

func NewDescription(in string) (Description, error) {
	d := Description(in)
	return d, d.Validate()
}

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

	return errors.New(primitives.InvalidProjectStatusCause, "Invalid project status: %s", p)
}

type Project struct {
	Name        primitives.DisplayName `json:"name",stbl:"name"`
	Label       primitives.Label       `json:"label"stbl:"label"`
	Description Description            `json:"description",stbl:"description"`
	Status      ProjectStatus          `json:"status",stbl:"status"`

	// The active spec (this references a ProjectSpec)
	CurrentSpecID *database.ID `stbl:"current_spec_id"`
}

func (p *Project) Validate() error {
	if err := p.Name.Validate(); err != nil {
		return err
	}

	if err := p.Label.Validate(); err != nil {
		return err
	}

	if p.CurrentSpecID != nil {
		if err := p.CurrentSpecID.Validate(); err != nil {
			return err
		}

		t, err := p.CurrentSpecID.Type()
		if err != nil {
			return err
		}

		if t != primitives.ProjectSpecType {
			return errors.New(primitives.InvalidIDCause, "CurrentSpecID can only be a ProjectSpec")
		}
	}

	if err := p.Status.Validate(); err != nil {
		return err
	}

	if err := p.Description.Validate(); err != nil {
		return err
	}

	return nil
}

func (p *Project) GetType() types.Type {
	return primitives.ProjectType
}

func NewProject(name primitives.DisplayName, label primitives.Label, description Description) *Project {
	return &Project{
		Name:        name,
		Label:       label,
		Description: description,
		Status:      ProjectPending,
	}
}
