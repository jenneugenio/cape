package integration

import (
	"context"
	capepg "github.com/capeprivacy/cape/coordinator/db/postgres"
	"github.com/capeprivacy/cape/models"
	gm "github.com/onsi/gomega"
	"time"

	"testing"
)

func TestContributor(t *testing.T) {
	gm.RegisterTestingT(t)
	ctx := context.TODO()

	testDB, err := CreateTestDB()
	gm.Expect(err).To(gm.BeNil())
	err = testDB.Setup(ctx)

	gm.Expect(err).To(gm.BeNil())
	defer testDB.Teardown(ctx) // nolint: errcheck

	cape := capepg.New(testDB.Pool)

	user := models.User{
		ID:        models.NewID(),
		Email:     "me@cape.com",
		Name:      "Me Me",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = cape.Users().Create(ctx, user)
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can't create a contributor if the project doesn't exist", func(t *testing.T) {
		_, err := cape.Contributors().Add(ctx, "my-contributor", "me@google.com")
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
