package encrypt

import (
	"context"
	"crypto/rand"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/manifoldco/go-base64"
)

var _ db.ProjectsDB = &projectEncrypt{}

type projectEncrypt struct {
	db    db.ProjectsDB
	codec crypto.EncryptionCodec
}

func (p *projectEncrypt) GetByID(ctx context.Context, id string) (*models.Project, error) {
	return p.db.GetByID(ctx, id)
}

func (p *projectEncrypt) Get(ctx context.Context, label models.Label) (*models.Project, error) {
	return p.db.Get(ctx, label)
}

func (p *projectEncrypt) Create(ctx context.Context, project models.Project) error {
	return p.db.Create(ctx, project)
}

func (p *projectEncrypt) Update(ctx context.Context, project models.Project) error {
	return p.db.Update(ctx, project)
}

func (p *projectEncrypt) CreateProjectSpec(ctx context.Context, spec models.Policy) error {
	for _, transform := range spec.Transformations {
		for key, arg := range transform.Args {
			sec, ok := arg.(models.SecretArg)
			if ok {
				// generate random bytes for secret value
				b := make([]byte, auth.SecretLength)
				_, err := rand.Read(b)
				if err != nil {
					return err
				}
				secValue := base64.New(b)

				// encrypt those random bytes for insertion into
				// database
				secValue, err = p.codec.Encrypt(ctx, secValue)
				if err != nil {
					return err
				}
				sec.Value = secValue

				transform.Args[key] = sec
			}
		}
	}
	return p.db.CreateProjectSpec(ctx, spec)
}

func (p *projectEncrypt) GetProjectSpec(ctx context.Context, id string) (*models.Policy, error) {
	spec, err := p.db.GetProjectSpec(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, transform := range spec.Transformations {
		for key, arg := range transform.Args {
			sec, ok := arg.(models.SecretArg)
			if ok {
				secValue, err := p.codec.Decrypt(ctx, sec.Value)
				if err != nil {
					return nil, err
				}

				sec.Value = secValue
				transform.Args[key] = sec
			}
		}
	}

	return spec, nil
}

func (p *projectEncrypt) List(ctx context.Context) ([]models.Project, error) {
	return p.db.List(ctx)
}

func (p *projectEncrypt) ListByStatus(ctx context.Context, status models.ProjectStatus) ([]models.Project, error) {
	return p.db.ListByStatus(ctx, status)
}

func (p *projectEncrypt) CreateSuggestion(ctx context.Context, suggestion models.Suggestion) error {
	return p.db.CreateSuggestion(ctx, suggestion)
}

func (p *projectEncrypt) GetSuggestions(ctx context.Context, projectLabel models.Label) ([]models.Suggestion, error) {
	return p.db.GetSuggestions(ctx, projectLabel)
}

func (p *projectEncrypt) GetSuggestion(ctx context.Context, id string) (*models.Suggestion, error) {
	return p.db.GetSuggestion(ctx, id)
}

func (p *projectEncrypt) UpdateSuggestion(ctx context.Context, suggestion models.Suggestion) error {
	return p.db.UpdateSuggestion(ctx, suggestion)
}
