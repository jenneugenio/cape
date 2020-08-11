package models

import (
	"encoding/json"
	"testing"

	gm "github.com/onsi/gomega"
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
			cause string
		}{
			"missing scheme": {
				in:    "s",
				cause: "invalid db url",
			},
			"wrong scheme": {
				in:    "ftp://my.coordinator.com",
				cause: "invalid db url",
			},
			"missing host": {
				in:    "postgres://",
				cause: "invalid db url",
			},
			"invalid host": {
				in:    "postgres://1323",
				cause: "invalid db url",
			},
			"missing path": {
				in:    "postgres://hello",
				cause: "invalid db url",
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				_, err := NewDBURL(test.in)
				gm.Expect(err).ToNot(gm.BeNil())
				gm.Expect(err.Error()).To(gm.Equal(test.cause))
			})
		}
	})
}
