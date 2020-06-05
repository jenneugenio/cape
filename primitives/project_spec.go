package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"sigs.k8s.io/yaml"
)

type ProjectSpecFile struct {
	Sources []Label `json:"sources"`
	Policy  []*Rule `json:"policy"`
}

func ParseProjectSpecFile(data []byte) (*ProjectSpecFile, error) {
	var spec ProjectSpecFile
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

type ProjectSpec struct {
	*database.Primitive
	ProjectID database.ID
	ParentID  *database.ID
	SourceIDs []database.ID
	Policy    []*Rule `json:"policy"`
}

func (p *ProjectSpec) GetType() types.Type {
	return ProjectSpecType
}

func (p *ProjectSpec) Validate() error {
	if err := p.ProjectID.Validate(); err != nil {
		return err
	}

	t, err := p.ProjectID.Type()
	if err != nil {
		return err
	}
	if t != ProjectType {
		return errors.New(InvalidIDCause, "ProjectID can only be a Project")
	}

	if p.ParentID != nil {
		if err := p.ParentID.Validate(); err != nil {
			return err
		}

		t, err := p.ParentID.Type()
		if err != nil {
			return err
		}

		if t != ProjectSpecType {
			return errors.New(InvalidIDCause, "ParentID can only be a ProjectSpec")
		}
	}

	for _, s := range p.SourceIDs {
		if err := s.Validate(); err != nil {
			return err
		}

		t, err := s.Type()
		if err != nil {
			return err
		}

		if t != SourcePrimitiveType {
			return errors.New(InvalidIDCause, "SourceIDs can only contain Source IDs")
		}
	}

	if len(p.Policy) == 0 {
		return errors.New(InvalidProjectSpecCause, "ProjectSpecs must define at least one policy")
	}

	for _, policy := range p.Policy {
		err := policy.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewProjectSpec(
	projectID database.ID,
	parent *database.ID,
	sources []database.ID,
	policy []*Rule,
) (*ProjectSpec, error) {
	p, err := database.NewPrimitive(ProjectSpecType)
	if err != nil {
		return nil, err
	}

	spec := &ProjectSpec{
		Primitive: p,
		ProjectID: projectID,
		ParentID:  parent,
		SourceIDs: sources,
		Policy:    policy,
	}

	if err := spec.Validate(); err != nil {
		return nil, err
	}

	return spec, nil
}
