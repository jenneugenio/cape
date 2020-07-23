package models

import (
	"time"

	"sigs.k8s.io/yaml"
)

type ProjectSpecFile struct {
	Transformations []NamedTransformation `json:"transformations"`
	Policy          []*Rule               `json:"policy"`
}

func ParseProjectSpecFile(data []byte) (*ProjectSpecFile, error) {
	var spec ProjectSpecFile
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

type ProjectSpec struct {
	ID              string                `json:"id"`
	ProjectID       string                `json:"project_id"`
	ParentID        *string               `json:"parent_id"`
	Transformations []NamedTransformation `json:"transformations"`
	Policy          []*Rule               `json:"policy"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

func (p *ProjectSpec) Validate() error {
	return nil
}

func NewProjectSpec(
	projectID string,
	parent *string,
	policy []*Rule,
	named []NamedTransformation,
) ProjectSpec {
	return ProjectSpec{
		ID:              NewID(),
		CreatedAt:       now(),
		ProjectID:       projectID,
		ParentID:        parent,
		Policy:          policy,
		Transformations: named,
	}
}
