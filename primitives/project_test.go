package primitives

import (
	"testing"

	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
)

func TestProject(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("Can create a new project with valid arguments", func(t *testing.T) {
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", "cool project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p).ToNot(gm.BeNil())
	})

	t.Run("New projects start in a pending state", func(t *testing.T) {
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", "cool project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.Status).To(gm.Equal(ProjectPending))
	})

	t.Run("New projects get an ID", func(t *testing.T) {
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", "cool project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.ID).ToNot(gm.BeNil())
	})

	t.Run("New projects returns the appropriate type", func(t *testing.T) {
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", "cool project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(p.GetType()).To(gm.Equal(ProjectType))
	})

	t.Run("Cannot set a fake status", func(t *testing.T) {
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", "cool project")
		gm.Expect(err).To(gm.BeNil())

		p.Status = ProjectStatus(":O")
		err = p.Validate()
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectStatusCause)).To(gm.BeTrue())
	})

	t.Run("Cannot use a huge description", func(t *testing.T) {
		myDesc := ""
		for i := 0; i <= maxDescriptionSize+1; i++ {
			myDesc += "x"
		}

		desc := Description(myDesc)
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", desc)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectDescriptionCause)).To(gm.BeTrue())
		gm.Expect(p).To(gm.BeNil())
	})

	t.Run("Validation fails if the `current` value is not a project spec", func(t *testing.T) {
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", "cool project")
		gm.Expect(err).To(gm.BeNil())

		email, err := NewEmail("hacker@cape.com")
		gm.Expect(err).To(gm.BeNil())

		creds, err := GenerateCredentials()
		gm.Expect(err).To(gm.BeNil())

		user, err := NewUser("Science McGee", email, creds)
		gm.Expect(err).To(gm.BeNil())

		p.CurrentSpecID = &user.ID
		err = p.Validate()
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidIDCause)).To(gm.BeTrue())
	})

	t.Run("Passes validation when `current` is a project spec", func(t *testing.T) {
		p, err := NewProject("Credit Card Recommendations", "credit-card-recommendations", "cool project")
		gm.Expect(err).To(gm.BeNil())

		policy, err := samplePolicySpec()
		gm.Expect(err).To(gm.BeNil())

		projectSpec, err := NewProjectSpec(p.ID, nil, policy.Rules)
		gm.Expect(err).To(gm.BeNil())

		p.CurrentSpecID = &projectSpec.ID
		err = p.Validate()
		gm.Expect(err).To(gm.BeNil())
	})
}

func TestInvalidProjects(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		Name          string
		ProjectName   DisplayName
		Label         Label
		Description   Description
		ExpectedCause errors.Cause
	}{
		{
			Name:          "Does not allow short project names",
			ProjectName:   "a",
			Label:         "my-project",
			Description:   "weep woop",
			ExpectedCause: InvalidProjectNameCause,
		},

		{
			Name:          "Does not allow an invalid label",
			ProjectName:   "Big Science Brain",
			Label:         "^_________^",
			Description:   "Lots of science happens in this project",
			ExpectedCause: InvalidLabelCause,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p, err := NewProject(test.ProjectName, test.Label, test.Description)
			gm.Expect(p).To(gm.BeNil())
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(errors.CausedBy(err, test.ExpectedCause)).To(gm.BeTrue())
		})
	}
}

func TestValidProjectNames(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		Name        string
		ProjectName DisplayName
		Label       Label
		Description Description
	}{
		{
			Name:        "Regular name",
			ProjectName: "Grade 5 Science Project",
			Label:       "my-project",
			Description: "Fun project where the science goes down",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p, err := NewProject(test.ProjectName, test.Label, test.Description)
			gm.Expect(err).To(gm.BeNil())
			gm.Expect(p).ToNot(gm.BeNil())
		})
	}
}

func TestProjectStatus(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Cannot make a fake status", func(t *testing.T) {
		p := ProjectStatus("haha not real!")
		err := p.Validate()
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectStatusCause)).To(gm.BeTrue())
	})
}

func TestProjectDescription(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("Cannot max a huge fake status", func(t *testing.T) {
		myDesc := ""
		for i := 0; i <= maxDescriptionSize+1; i++ {
			myDesc += "x"
		}

		_, err := NewDescription(myDesc)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidProjectDescriptionCause)).To(gm.BeTrue())
	})
}
