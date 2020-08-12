package capepg

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	gm "github.com/onsi/gomega"
	"reflect"
	"testing"
)

type projectTestPool struct {
	err error
}

type ProjectRow struct {
	err error
}

func (p *ProjectRow) Scan(dest ...interface{}) error {
	if p.err != nil {
		return p.err
	}

	project := models.Project{
		ID:          "myproject",
		Label:       "myproject",
		Name:        "myproject",
		Description: "its my project",
	}

	for _, item := range dest {
		orig := reflect.ValueOf(item)
		replacement := project
		reflect.Indirect(orig).Set(reflect.ValueOf(replacement))
	}

	return nil
}

func (p projectTestPool) Exec(_ context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return nil, p.err
}

func (p projectTestPool) Query(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	panic("implement me")
}

func (p projectTestPool) QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row {
	return &ProjectRow{p.err}
}

func TestProject(t *testing.T) {
	gm.RegisterTestingT(t)

	p := pgProject{
		pool: &projectTestPool{},
	}

	t.Run("Can get by ID", func(t *testing.T) {
		project, err := p.GetByID(context.TODO(), "abc123")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(project.ID).To(gm.Equal("myproject"))

		t.Run("returns a special error if no project is found", func(t *testing.T) {
			p := pgProject{
				pool: &projectTestPool{
					err: pgx.ErrNoRows,
				},
			}

			_, err := p.GetByID(context.TODO(), "abc123")
			gm.Expect(err).To(gm.Equal(db.ErrCannotFindProject))
		})
	})

	t.Run("Can get by label", func(t *testing.T) {
		project, err := p.Get(context.TODO(), "abc123")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(project.ID).To(gm.Equal("myproject"))

		t.Run("returns a special error if no project is found", func(t *testing.T) {
			p := pgProject{
				pool: &projectTestPool{
					err: pgx.ErrNoRows,
				},
			}

			_, err := p.Get(context.TODO(), "abc123")
			gm.Expect(err).To(gm.Equal(db.ErrCannotFindProject))
		})
	})

	t.Run("Can create", func(t *testing.T) {
		err := p.Create(context.TODO(), models.Project{})
		gm.Expect(err).To(gm.BeNil())

		t.Run("Cannot create duplicate projects", func(t *testing.T) {
			p := pgProject{
				pool: &projectTestPool{
					err: fmt.Errorf("can't do it"),
				},
			}
			err := p.Create(context.TODO(), models.Project{})
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(err.Error()).To(gm.Equal("entity already exists"))
		})
	})

	t.Run("Can update", func(t *testing.T) {
		err := p.Update(context.TODO(), models.Project{})
		gm.Expect(err).To(gm.BeNil())
	})
}
