// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
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
		email, err := primitives.NewEmail("deleted-service@connector-cape.com")
		gm.Expect(err).To(gm.BeNil())

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
		email, err := primitives.NewEmail("service:email@connector-cape.com")
		gm.Expect(err).To(gm.BeNil())

		s, err := createServicePrimitive(email)
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

		s, err := createServicePrimitive(email)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		email, err = primitives.NewEmail("fresh-email@bomb.com")
		gm.Expect(err).To(gm.BeNil())
		s, err = createServicePrimitive(email)
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(service).To(gm.BeNil())
	})

	t.Run("can create data connector service", func(t *testing.T) {
		email, err := primitives.NewEmail("dc@cape.com")
		gm.Expect(err).To(gm.BeNil())

		url, err := primitives.NewURL("https://cape.com")
		gm.Expect(err).To(gm.BeNil())

		s, err := primitives.NewService(email, primitives.DataConnectorServiceType, url)
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(service.Type).To(gm.Equal(s.Type))
		gm.Expect(service.Endpoint).To(gm.Equal(s.Endpoint))

		connectorRole, err := client.GetRoleByLabel(ctx, primitives.DataConnectorRole)
		gm.Expect(err).To(gm.BeNil())

		members, err := client.GetMembersRole(ctx, connectorRole.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(members[0].GetEmail()).To(gm.Equal(service.Email))
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

	globalRole, err := client.GetRoleByLabel(ctx, primitives.GlobalRole)
	gm.Expect(err).To(gm.BeNil())

	cRole, err := client.CreateRole(ctx, "connector", nil)
	gm.Expect(err).To(gm.BeNil())

	dsRole, err := client.CreateRole(ctx, "ds-role", nil)
	gm.Expect(err).To(gm.BeNil())

	emails := []string{"connector1@email.com", "connector2@email.com", "connector3@email.com"}
	services := make([]*coordinator.ServiceResponse, 3)
	for i, email := range emails {
		e, err := primitives.NewEmail(email)
		gm.Expect(err).To(gm.BeNil())

		s, err := createServicePrimitive(e)
		gm.Expect(err).To(gm.BeNil())

		service, err := client.CreateService(ctx, s)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AssignRole(ctx, service.ID, cRole.ID)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.AssignRole(ctx, service.ID, dsRole.ID)
		gm.Expect(err).To(gm.BeNil())

		services[i] = &coordinator.ServiceResponse{
			Service: service,
			Roles:   []*primitives.Role{globalRole, cRole, dsRole},
		}
	}

	otherServices, err := client.ListServices(ctx)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(otherServices).To(gm.ContainElements(services))
}

func createServicePrimitive(email primitives.Email) (*primitives.Service, error) {
	typ, err := primitives.NewServiceType("user")
	if err != nil {
		return nil, err
	}

	service, err := primitives.NewService(email, typ, nil)
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

	s, err := createServicePrimitive(email)
	gm.Expect(err).To(gm.BeNil())

	service, err := client.CreateService(ctx, s)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(service.Email).To(gm.Equal(email))

	token, err := client.NewToken(ctx, service)
	gm.Expect(err).To(gm.BeNil())

	serviceClient, err := h.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = serviceClient.TokenLogin(ctx, token)
	gm.Expect(err).To(gm.BeNil())

	sources, err := serviceClient.ListSources(ctx)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(len(sources)).To(gm.Equal(0))
}
