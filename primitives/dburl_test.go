package primitives

import (
	"encoding/json"
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestDBURL(t *testing.T) {
	gm.RegisterTestingT(t)

	good := "postgres://u:p@host.com:5432/hello"
	t.Run("parses a valid addr", func(t *testing.T) {
		out, err := NewDBURL(good)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(out.String()).To(gm.Equal(good))
	})

	t.Run("marshal to json", func(t *testing.T) {
		out, err := NewDBURL(good)
		gm.Expect(err).To(gm.BeNil())

		result, err := json.Marshal(out)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(string(result)).To(gm.Equal("\"" + good + "\""))
	})

	t.Run("unmarshal from json", func(t *testing.T) {
		out, err := NewDBURL(good)
		gm.Expect(err).To(gm.BeNil())

		result, err := json.Marshal(out)
		gm.Expect(err).To(gm.BeNil())

		new := &DBURL{}
		err = json.Unmarshal(result, new)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(new.String()).To(gm.Equal(good))
		gm.Expect(new.Validate()).To(gm.BeNil())
	})

	t.Run("catches bad errors", func(t *testing.T) {
		tests := map[string]struct {
			in    string
			cause errors.Cause
		}{
			"missing scheme": {
				in:    "s",
				cause: InvalidDBURLCause,
			},
			"wrong scheme": {
				in:    "ftp://my.coordinator.com",
				cause: InvalidDBURLCause,
			},
			"missing host": {
				in:    "postgres://",
				cause: InvalidDBURLCause,
			},
			"invalid host": {
				in:    "postgres://1323",
				cause: InvalidDBURLCause,
			},
			"missing path": {
				in:    "postgres://hello",
				cause: InvalidDBURLCause,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				_, err := NewDBURL(test.in)
				gm.Expect(errors.FromCause(err, test.cause)).To(gm.BeTrue())
			})
		}
	})
}
