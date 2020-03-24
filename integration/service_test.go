// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
)

func TestServices(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = client.Login(ctx, tc.User.Email, tc.UserPassword)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create service", func(t *testing.T) {
		email := "service@connector-cape.com"
		s, err := createServicePrimitive(email)
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(service.Email).To(gm.Equal(email))

		otherService, err := client.GetService(ctx, service.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(service.Email).To(gm.Equal(otherService.Email))
		gm.Expect(service.ID).To(gm.Equal(otherService.ID))
	})

	t.Run("delete service", func(t *testing.T) {
		email := "deleted-service@connector-cape.com"
		s, err := createServicePrimitive(email)
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		err = client.DeleteService(ctx, service.ID)
		gm.Expect(err).To(gm.BeNil())

		otherService, err := client.GetService(ctx, service.ID)
		gm.Expect(err).NotTo(gm.BeNil())
		gm.Expect(otherService).To(gm.BeNil())
	})

	t.Run("get service by email", func(t *testing.T) {
		email := "email@connector-cape.com"
		s, err := createServicePrimitive(email)
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		otherService, err := client.GetServiceByEmail(ctx, email)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(service.Email).To(gm.Equal(otherService.Email))
		gm.Expect(service.ID).To(gm.Equal(otherService.ID))
	})
}

func TestListServices(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer tc.Teardown(ctx) // nolint: errcheck

	client, err := tc.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = client.Login(ctx, tc.User.Email, tc.UserPassword)
	gm.Expect(err).To(gm.BeNil())

	emails := []string{"connector1@email.com", "connector2@email.com", "connector3@email.com"}
	services := make([]*primitives.Service, 3)
	for i, email := range emails {
		s, err := createServicePrimitive(email)
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		services[i] = service
	}

	otherServices, err := client.ListServices(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(otherServices).To(gm.ContainElements(services))
}

func createServicePrimitive(email string) (*primitives.Service, error) {
	creds, err := auth.NewCredentials([]byte("connectorsecretsarecool"), nil)
	if err != nil {
		return nil, err
	}

	service, err := primitives.NewService(email, creds.Package())
	if err != nil {
		return nil, err
	}

	return service, nil
}
