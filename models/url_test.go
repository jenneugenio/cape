package models

import (
	"encoding/json"
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestNewURL(t *testing.T) {
	gm.RegisterTestingT(t)

	u := "https://my.coordinator.com"
	t.Run("parses a valid url", func(t *testing.T) {
		out, err := NewURL(u)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(out.String()).To(gm.Equal(u))
	})

	t.Run("marshal to json", func(t *testing.T) {
		out, err := NewURL(u)
		gm.Expect(err).To(gm.BeNil())

		result, err := json.Marshal(out)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(string(result)).To(gm.Equal("\"" + u + "\""))
	})

	t.Run("unmarshal from json", func(t *testing.T) {
		out, err := NewURL(u)
		gm.Expect(err).To(gm.BeNil())

		result, err := json.Marshal(out)
		gm.Expect(err).To(gm.BeNil())

		new := &URL{}
		err = json.Unmarshal(result, new)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(new.String()).To(gm.Equal(u))
		gm.Expect(new.Validate()).To(gm.BeNil())
	})

	t.Run("catches bad errors", func(t *testing.T) {
		tests := map[string]struct {
			in    string
			cause errors.Cause
		}{
			"missing scheme": {
				in:    "s",
				cause: InvalidURLCause,
			},
			"wrong scheme": {
				in:    "ftp://my.coordinator.com",
				cause: InvalidURLCause,
			},
			"missing host": {
				in:    "https://",
				cause: InvalidURLCause,
			},
			"invalid host": {
				in:    "postgres://1323",
				cause: InvalidURLCause,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				_, err := NewURL(test.in)
				gm.Expect(errors.FromCause(err, test.cause)).To(gm.BeTrue())
			})
		}
	})
}
