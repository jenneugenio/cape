// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/controller/harness"
	"github.com/dropoutlabs/cape/primitives"
)

func TestServices(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create service", func(t *testing.T) {
		email, err := primitives.NewEmail("service@connector-cape.com")
		gm.Expect(err).To(gm.BeNil())

		s, err := createServicePrimitive(email, []byte("randompassword"))
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
		email, err := primitives.NewEmail("deleted-service@connector-cape.com")
		gm.Expect(err).To(gm.BeNil())

		s, err := createServicePrimitive(email, []byte("randompassword"))
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
		email, err := primitives.NewEmail("service:email@connector-cape.com")
		gm.Expect(err).To(gm.BeNil())

		s, err := createServicePrimitive(email, []byte("randompassword"))
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		otherService, err := client.GetServiceByEmail(ctx, email)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(service.Email).To(gm.Equal(otherService.Email))
		gm.Expect(service.ID).To(gm.Equal(otherService.ID))
	})

	t.Run("cannot create multiple services with same email", func(t *testing.T) {
		email, err := primitives.NewEmail("fresh-email@bomb.com")
		gm.Expect(err).To(gm.BeNil())

		s, err := createServicePrimitive(email, []byte("randompassword"))
		gm.Expect(err).To(gm.BeNil())

		_, err = client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		email, err = primitives.NewEmail("fresh-email@bomb.com")
		gm.Expect(err).To(gm.BeNil())
		s, err = createServicePrimitive(email, []byte("randompassword"))
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(service).To(gm.BeNil())
	})
}

func TestListServices(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	cRole, err := client.CreateRole(ctx, "connector", nil)
	gm.Expect(err).To(gm.BeNil())

	dsRole, err := client.CreateRole(ctx, "ds-role", nil)
	gm.Expect(err).To(gm.BeNil())

	emails := []string{"connector1@email.com", "connector2@email.com", "connector3@email.com"}
	services := make([]*controller.ServiceResponse, 3)
	for i, email := range emails {
		e, err := primitives.NewEmail(email)
		gm.Expect(err).To(gm.BeNil())

		s, err := createServicePrimitive(e, []byte("randompassword"))
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AssignRole(ctx, service.ID, cRole.ID)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AssignRole(ctx, service.ID, dsRole.ID)
		gm.Expect(err).To(gm.BeNil())

		services[i] = &controller.ServiceResponse{
			Service: service,
			Roles:   []*primitives.Role{cRole, dsRole},
		}
	}

	otherServices, err := client.ListServices(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(otherServices).To(gm.ContainElements(services))
}

func createServicePrimitive(email primitives.Email, secret []byte) (*primitives.Service, error) {
	creds, err := auth.NewCredentials(secret, nil)
	if err != nil {
		return nil, err
	}

	typ, err := primitives.NewServiceType("user")
	if err != nil {
		return nil, err
	}

	service, err := primitives.NewService(email, typ, creds.Package())
	if err != nil {
		return nil, err
	}

	return service, nil
}

func TestServiceLogin(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	email, err := primitives.NewEmail("service:service@connector-cape.com")
	gm.Expect(err).To(gm.BeNil())

	url, err := h.URL()
	gm.Expect(err).To(gm.BeNil())

	token, err := auth.NewAPIToken(email, url)
	gm.Expect(err).To(gm.BeNil())

	s, err := createServicePrimitive(email, token.Secret)
	gm.Expect(err).To(gm.BeNil())

	service, err := client.CreateService(ctx, s)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(service.Email).To(gm.Equal(email))

	serviceClient, err := h.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = serviceClient.Login(ctx, token.Email, token.Secret)
	gm.Expect(err).To(gm.BeNil())

	sources, err := serviceClient.ListSources(ctx)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(len(sources)).To(gm.Equal(0))
}
