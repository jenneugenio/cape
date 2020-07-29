package encrypt

import (
	"context"
	"reflect"
	"testing"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/models"
	"github.com/manifoldco/go-base64"

	gm "github.com/onsi/gomega"
)

var SecretProjectSpec = models.Policy{
	ID: "ididid",
	Transformations: []*models.NamedTransformation{
		{
			Name: "TestTransform",
			Type: "secret",
			Args: map[string]interface{}{
				"key": &models.SecretArg{
					Type: "secret",
					Name: "my-key",
				},
				"nonSecret": 10,
			},
		},
	},
}

func TestProjectSpecCreate(t *testing.T) {
	tests := []struct {
		spec    models.Policy
		wantErr error
		err     error
	}{
		{
			spec:    SecretProjectSpec,
			wantErr: nil,
			err:     nil,
		},
		{
			spec:    SecretProjectSpec,
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
	}

	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	pgProjects := &testPgProjects{}
	for i, test := range tests {
		projectDB := projectEncrypt{
			db:    pgProjects,
			codec: codec,
		}

		pgProjects.err = test.err

		gotErr := projectDB.CreateProjectSpec(context.TODO(), test.spec)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}

		if !reflect.DeepEqual(SecretProjectSpec, pgProjects.receivedProjectSpec) {
			t.Errorf("project spec not encrypted: got %v", pgProjects.receivedProjectSpec)
		}
	}
}

func TestProjectSpecGet(t *testing.T) {
	gm.RegisterTestingT(t)

	secret := base64.New([]byte("secretsecret"))

	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	encryptedSecret, _ := codec.Encrypt(context.TODO(), secret)
	encryptedSpec := models.Policy{
		ID: "ididid",
		Transformations: []*models.NamedTransformation{
			{
				Name: "TestTransform",
				Type: "secret",
				Args: map[string]interface{}{
					"key": models.SecretArg{
						Type:  "secret",
						Name:  "my-key",
						Value: encryptedSecret,
					},
				},
			},
		},
	}

	secretSpec := models.Policy{
		ID: "ididid",
		Transformations: []*models.NamedTransformation{
			{
				Name: "TestTransform",
				Type: "secret",
				Args: map[string]interface{}{
					"key": models.SecretArg{
						Type:  "secret",
						Name:  "my-key",
						Value: secret,
					},
				},
			},
		},
	}

	tests := []struct {
		spec     models.Policy
		wantSpec *models.Policy
		wantErr  error
		err      error
	}{
		{
			spec:     encryptedSpec,
			wantSpec: &secretSpec,
			wantErr:  nil,
			err:      nil,
		},
		{
			spec:    encryptedSpec,
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
	}

	pgProjects := &testPgProjects{}
	for i, test := range tests {
		projectDB := projectEncrypt{
			db:    pgProjects,
			codec: codec,
		}

		pgProjects.returnProjectSpec = test.spec
		pgProjects.err = test.err

		gotSpec, gotErr := projectDB.GetProjectSpec(context.TODO(), test.spec.ID)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Get() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}

		gm.Expect(gotSpec).To(gm.Equal(test.wantSpec))
	}
}

type testPgProjects struct {
	returnProjectSpec   models.Policy
	receivedProjectSpec models.Policy
	err                 error
}

// only testing the below two for now, rest can remain unimplemented

func (t *testPgProjects) CreateProjectSpec(ctx context.Context, spec models.Policy) error {
	t.receivedProjectSpec = spec
	return t.err
}

func (t *testPgProjects) GetProjectSpec(ctx context.Context, spec string) (*models.Policy, error) {
	if t.err != nil {
		return nil, t.err
	}

	return &t.returnProjectSpec, nil
}

func (t *testPgProjects) Get(_ context.Context, _ models.Label) (*models.Project, error) {
	panic("not implemented")
}

func (t *testPgProjects) GetByID(_ context.Context, _ string) (*models.Project, error) {
	panic("not implemented")
}

func (t *testPgProjects) Create(_ context.Context, _ models.Project) error {
	panic("not implemented")
}

func (t *testPgProjects) Update(_ context.Context, _ models.Project) error {
	panic("not implemented")
}

func (t *testPgProjects) List(_ context.Context) ([]models.Project, error) {
	panic("not implemented")
}

func (t *testPgProjects) ListByStatus(_ context.Context, _ models.ProjectStatus) ([]models.Project, error) {
	panic("not implemented")
}

func (t *testPgProjects) CreateSuggestion(_ context.Context, _ models.Suggestion) error {
	panic("not implemented")
}
