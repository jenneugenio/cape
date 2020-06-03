package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
	"io/ioutil"
	"testing"
)

func samplePolicySpec() ([]*PolicySpec, error) {
	specBytes, err := ioutil.ReadFile("./testdata/policy.yaml")
	if err != nil {
		return nil, err
	}

	spec, err := ParsePolicySpec(specBytes)
	if err != nil {
		return nil, err
	}
	return []*PolicySpec{spec}, nil
}

func TestProjectSpec(t *testing.T) {
	gm.RegisterTestingT(t)

	projectID, err := database.GenerateID(ProjectType)
	gm.Expect(err).To(gm.BeNil())

	sourceID, err := database.GenerateID(SourcePrimitiveType)
	gm.Expect(err).To(gm.BeNil())

	policy, err := samplePolicySpec()
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can create a project spec without a parent", func(t *testing.T) {
		ps, err := NewProjectSpec(projectID, nil, []database.ID{sourceID}, policy)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(ps).ToNot(gm.BeNil())
		gm.Expect(ps.ID).ToNot(gm.BeNil())
	})

	t.Run("Can create a project spec with a parent", func(t *testing.T) {
		parent, err := NewProjectSpec(projectID, nil, []database.ID{sourceID}, policy)
		gm.Expect(err).To(gm.BeNil())

		ps, err := NewProjectSpec(projectID, &parent.ID, []database.ID{sourceID}, policy)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(ps).ToNot(gm.BeNil())
		gm.Expect(ps.ID).ToNot(gm.BeNil())
		gm.Expect(ps.ParentID).To(gm.Equal(&parent.ID))
	})

	t.Run("Cannot create a project spec without a policy", func(t *testing.T) {
		ps, err := NewProjectSpec(projectID, nil, []database.ID{sourceID}, nil)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectSpecCause)).To(gm.BeTrue())
		gm.Expect(ps).To(gm.BeNil())
	})
}

func TestProjectSpecInvalidIDs(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		Name      string
		ProjectID types.Type
		Parent    types.Type
		Sources   []types.Type
	}{
		{
			Name:      "Invalid ProjectID",
			ProjectID: UserType,
			Parent:    ProjectSpecType,
			Sources:   []types.Type{SourcePrimitiveType},
		},

		{
			Name:      "Invalid ParentID",
			ProjectID: ProjectType,
			Parent:    UserType,
			Sources:   []types.Type{SourcePrimitiveType},
		},

		{
			Name:      "Invalid Single SourceID",
			ProjectID: ProjectType,
			Parent:    ProjectSpecType,
			Sources:   []types.Type{UserType},
		},

		{
			Name:      "One invalid SourceID amongst many",
			ProjectID: ProjectType,
			Parent:    ProjectSpecType,
			Sources:   []types.Type{SourcePrimitiveType, SourcePrimitiveType, UserType},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			projectID, err := database.GenerateID(test.ProjectID)
			gm.Expect(err).To(gm.BeNil())

			ParentID, err := database.GenerateID(test.Parent)
			gm.Expect(err).To(gm.BeNil())

			sourceIDs := make([]database.ID, len(test.Sources))
			for i, s := range test.Sources {
				sid, err := database.GenerateID(s)
				gm.Expect(err).To(gm.BeNil())

				sourceIDs[i] = sid
			}

			policies, err := samplePolicySpec()
			gm.Expect(err).To(gm.BeNil())

			projectSpec, err := NewProjectSpec(projectID, &ParentID, sourceIDs, policies)
			gm.Expect(projectSpec).To(gm.BeNil())
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(errors.CausedBy(err, InvalidIDCause)).To(gm.BeTrue())
		})
	}
}
