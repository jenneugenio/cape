// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/coordinator/client"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
)

func TestIdentities(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err).To(gm.BeNil())

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	services, err := createServices(ctx, client, 3)
	gm.Expect(err).To(gm.BeNil())

	users, err := createUsers(ctx, client, 2)
	gm.Expect(err).To(gm.BeNil())

	var emails []primitives.Email
	for _, service := range services {
		emails = append(emails, service.Email)
	}

	for _, user := range users {
		emails = append(emails, user.Email)
	}

	identities, err := client.GetIdentities(ctx, emails)
	gm.Expect(err).To(gm.BeNil())

	var otherEmails []primitives.Email
	for _, identity := range identities {
		otherEmails = append(otherEmails, identity.GetEmail())
	}

	gm.Expect(otherEmails).To(gm.ContainElements(emails))
}

func createServices(ctx context.Context, client *client.Client, numServices int) ([]*primitives.Service, error) {
	services := make([]*primitives.Service, numServices)
	for i := 0; i < numServices; i++ {
		email, err := primitives.NewEmail(fmt.Sprintf("service:email%d@email.com", i))
		if err != nil {
			return nil, err
		}

		service, err := createServicePrimitive(email)
		if err != nil {
			return nil, err
		}

		_, err = client.CreateService(ctx, service)
		if err != nil {
			return nil, err
		}

		services[i] = service
	}

	return services, nil
}

func createUsers(ctx context.Context, client *client.Client, numUsers int) ([]*primitives.User, error) {
	users := make([]*primitives.User, numUsers)
	for i := 0; i < numUsers; i++ {
		email, err := primitives.NewEmail(fmt.Sprintf("email%d@email.com", i))
		if err != nil {
			return nil, err
		}

		name, err := primitives.NewName(fmt.Sprintf("Hi%d Hello", i))
		if err != nil {
			return nil, err
		}

		user, _, err := client.CreateUser(ctx, name, email)
		if err != nil {
			return nil, err
		}

		users[i] = user
	}

	return users, nil
}
