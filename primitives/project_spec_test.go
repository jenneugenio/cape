package primitives

import (
	"io/ioutil"
	"testing"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
)

func samplePolicySpec() (*PolicySpec, error) {
	specBytes, err := ioutil.ReadFile("./testdata/policy.yaml")
	if err != nil {
		return nil, err
	}

	spec, err := ParsePolicySpec(specBytes)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

func TestProjectSpec(t *testing.T) {
	gm.RegisterTestingT(t)

	projectID, err := database.GenerateID(ProjectType)
	gm.Expect(err).To(gm.BeNil())

	policy, err := samplePolicySpec()
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can create a project spec without a parent", func(t *testing.T) {
		ps, err := NewProjectSpec(projectID, nil, policy.Rules)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(ps).ToNot(gm.BeNil())
		gm.Expect(ps.ID).ToNot(gm.BeNil())
	})

	t.Run("Can create a project spec with a parent", func(t *testing.T) {
		parent, err := NewProjectSpec(projectID, nil, policy.Rules)
		gm.Expect(err).To(gm.BeNil())

		ps, err := NewProjectSpec(projectID, &parent.ID, policy.Rules)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(ps).ToNot(gm.BeNil())
		gm.Expect(ps.ID).ToNot(gm.BeNil())
		gm.Expect(ps.ParentID).To(gm.Equal(&parent.ID))
	})

	t.Run("Cannot create a project spec without a policy", func(t *testing.T) {
		ps, err := NewProjectSpec(projectID, nil, nil)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectSpecCause)).To(gm.BeTrue())
		gm.Expect(ps).To(gm.BeNil())
	})

	t.Run("Can create a project spec from yaml", func(t *testing.T) {
		f, err := ioutil.ReadFile("./testdata/project_spec.yaml")
		gm.Expect(err).To(gm.BeNil())

		spec, err := ParseProjectSpecFile(f)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(spec).NotTo(gm.BeNil())
	})
}

func TestProjectSpecInvalidIDs(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		Name      string
		ProjectID types.Type
		Parent    types.Type
	}{
		{
			Name:      "Invalid ProjectID",
			ProjectID: UserType,
			Parent:    ProjectSpecType,
		},

		{
			Name:      "Invalid ParentID",
			ProjectID: ProjectType,
			Parent:    UserType,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			projectID, err := database.GenerateID(test.ProjectID)
			gm.Expect(err).To(gm.BeNil())

			ParentID, err := database.GenerateID(test.Parent)
			gm.Expect(err).To(gm.BeNil())

			policy, err := samplePolicySpec()
			gm.Expect(err).To(gm.BeNil())

			projectSpec, err := NewProjectSpec(projectID, &ParentID, policy.Rules)
			gm.Expect(projectSpec).To(gm.BeNil())
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(errors.CausedBy(err, InvalidIDCause)).To(gm.BeTrue())
		})
	}
}
