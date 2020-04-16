package primitives

import (
	"encoding/json"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/database"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestNewSource(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := map[string]struct {
		label       Label
		credentials string
		serviceID   *database.ID

		sourceType SourceType
		endpoint   string
		success    bool
		cause      errors.Cause
	}{
		"sets type properly": {
			label:       Label("hello"),
			credentials: "postgres://user:test@my.cool.website.com:5432/test",

			endpoint:   "postgres://my.cool.website.com:5432/test",
			sourceType: PostgresType,
			serviceID:  nil,
			success:    true,
		},
		"returns error if credentials is nil": {
			label:   Label("error-please"),
			success: false,
			cause:   InvalidSourceCause,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var u *DBURL
			var err error
			if test.credentials != "" {
				u, err = NewDBURL(test.credentials)
				gm.Expect(err).To(gm.BeNil())
			}

			source, err := NewSource(test.label, u, test.serviceID)
			if !test.success {
				gm.Expect(errors.FromCause(err, test.cause)).To(gm.BeTrue())
				return
			}

			gm.Expect(source.Label).To(gm.Equal(test.label))
			gm.Expect(source.Type).To(gm.Equal(test.sourceType))
			gm.Expect(source.ServiceID).To(gm.Equal(test.serviceID))
			gm.Expect(source.Credentials.String()).To(gm.Equal(test.credentials))
			gm.Expect(source.Endpoint.String()).To(gm.Equal(test.endpoint))
		})
	}

	u, err := NewDBURL("postgres://user:test@my.cool.com:5432/test")
	gm.Expect(err).To(gm.BeNil())

	t.Run("marshal to json", func(t *testing.T) {
		source, err := NewSource("helllllo", u, nil)
		gm.Expect(err).To(gm.BeNil())

		_, err = json.Marshal(source)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("unmarshal from json", func(t *testing.T) {
		source, err := NewSource("heyaaaa", u, nil)
		gm.Expect(err).To(gm.BeNil())

		b, err := json.Marshal(source)
		gm.Expect(err).To(gm.BeNil())

		out := &Source{}
		err = json.Unmarshal(b, out)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(source).To(gm.Equal(out))
	})
}
