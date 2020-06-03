package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

const maxDescriptionSize = 5000

type Description string

func (d Description) Validate() error {
	if len(d) > maxDescriptionSize {
		return errors.New(InvalidProjectDescriptionCause, "%d exceeds max description size of %d", len(d), maxDescriptionSize)
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

	return errors.New(InvalidProjectStatusCause, "Invalid project status: %s", p)
}

type Project struct {
	*database.Primitive
	Name        DisplayName   `json:"name"`
	Label       Label         `json:"label"`
	Description Description   `json:"description"`
	Status      ProjectStatus `json:"status"`

	// The active spec (this references a ProjectSpec)
	CurrentSpecID *database.ID
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

		if t != ProjectSpecType {
			return errors.New(InvalidIDCause, "CurrentSpecID can only be a ProjectSpec")
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
	return ProjectType
}

func NewProject(name DisplayName, label Label, description Description) (*Project, error) {
	p, err := database.NewPrimitive(ProjectType)
	if err != nil {
		return nil, err
	}

	project := &Project{
		Primitive:   p,
		Name:        name,
		Label:       label,
		Description: description,
		Status:      ProjectPending,
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}
	return project, nil
}
