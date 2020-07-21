package models

import "sigs.k8s.io/yaml"

type ProjectSpecFile struct {
	Policy []*Rule `json:"policy"`
}

func ParseProjectSpecFile(data []byte) (*ProjectSpecFile, error) {
	var spec ProjectSpecFile
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

type ProjectSpec struct {
	ID        string `json:"id"`
	ProjectID string
	ParentID  *string
	Policy    []*Rule `json:"policy"`
}

func (p *ProjectSpec) Validate() error {
	return nil
}

func NewProjectSpec(
	projectID string,
	parent *string,
	policy []*Rule,
) ProjectSpec {
	return ProjectSpec{
		ID:        NewID(),
		ProjectID: projectID,
		ParentID:  parent,
		Policy:    policy,
	}
}
