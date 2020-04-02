package primitives

import (
	"encoding/json"
	"net/url"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/database"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

func TestNewSource(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := map[string]struct {
		label       Label
		credentials string
		sourceType  SourceType
		serviceID   *database.ID

		success bool
		cause   errors.Cause
	}{
		"sets type properly": {
			label:       Label("hello"),
			credentials: "postgres://my.cool.website.com:5432/test",
			sourceType:  PostgresType,

			serviceID: nil,
			success:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(test.credentials)
			gm.Expect(err).To(gm.BeNil())

			source, err := NewSource(test.label, *u, test.serviceID)
			if !test.success {
				gm.Expect(errors.FromCause(err, test.cause)).To(gm.BeTrue())
				return
			}

			gm.Expect(source.Label).To(gm.Equal(test.label))
			gm.Expect(source.Type).To(gm.Equal(test.sourceType))
			gm.Expect(source.ServiceID).To(gm.Equal(test.serviceID))
		})
	}

	u, err := url.Parse("postgres://my.cool.com:5432/test")
	gm.Expect(err).To(gm.BeNil())

	t.Run("marshal to json", func(t *testing.T) {
		source, err := NewSource("helllllo", *u, nil)
		gm.Expect(err).To(gm.BeNil())

		_, err = json.Marshal(source)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("unmarshal from json", func(t *testing.T) {
		source, err := NewSource("heyaaaa", *u, nil)
		gm.Expect(err).To(gm.BeNil())

		b, err := json.Marshal(source)
		gm.Expect(err).To(gm.BeNil())

		out := &Source{}
		err = json.Unmarshal(b, out)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source).To(gm.Equal(out))
	})
}
