package models

import (
	"bytes"
	"encoding/json"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestMarshallNamedTransform(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("test json marshal", func(t *testing.T) {
		args := map[string]interface{}{
			"n": 100,
		}
		named := NamedTransformation{
			Name: "test",
			Type: "plusN",
			Args: args,
		}

		by, err := json.Marshal(named)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(string(by)).To(gm.Equal("{\"n\":100,\"name\":\"test\",\"type\":\"plusN\"}"))
	})

	t.Run("test json unmarshal", func(t *testing.T) {
		j := `
		{"n":100, "name": "test", "type": "plusN"}
		`

		var named NamedTransformation
		err := json.Unmarshal([]byte(j), &named)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(named.Name).To(gm.Equal("test"))
		gm.Expect(named.Type).To(gm.Equal("plusN"))
		gm.Expect(named.Args["n"]).To(gm.Equal(100.0))
	})

	t.Run("test json unmarshal with secret", func(t *testing.T) {
		j := `
		{
			"name": "test",
			"type": "tokenization",
			"key": {
				"name": "my-key",
				"type": "secret"
			},
			"nonSecret": 10
		}
		`

		var named NamedTransformation
		err := json.Unmarshal([]byte(j), &named)
		gm.Expect(err).To(gm.BeNil())

		sec := SecretArg{
			Type: "secret",
			Name: "my-key",
		}

		gm.Expect(named.Name).To(gm.Equal("test"))
		gm.Expect(named.Type).To(gm.Equal("tokenization"))
		gm.Expect(named.Args["key"]).To(gm.Equal(sec))
		gm.Expect(named.Args["nonSecret"]).To(gm.Equal(10.0))
	})

	t.Run("test gql marshal", func(t *testing.T) {
		args := map[string]interface{}{
			"n": 100,
		}
		named := NamedTransformation{
			Name: "test",
			Type: "plusN",
			Args: args,
		}

		buf := bytes.NewBuffer(nil)

		named.MarshalGQL(buf)

		gm.Expect(buf.String()).To(gm.Equal("{\"n\":100,\"name\":\"test\",\"type\":\"plusN\"}"))
	})

	t.Run("test gql unmarshal", func(t *testing.T) {
		namedMap := map[string]interface{}{
			"n":    100,
			"name": "test",
			"type": "plusN",
		}

		var named NamedTransformation
		err := named.UnmarshalGQL(namedMap)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(named.Name).To(gm.Equal("test"))
		gm.Expect(named.Type).To(gm.Equal("plusN"))
		gm.Expect(named.Args["n"]).To(gm.Equal(100))
	})

	t.Run("test gql unmarshal with secret", func(t *testing.T) {
		namedMap := map[string]interface{}{
			"key": map[string]interface{}{
				"name": "my-key",
				"type": "secret",
			},
			"nonSecret": 10,
			"name":      "test",
			"type":      "tokenization",
		}

		var named NamedTransformation
		err := named.UnmarshalGQL(namedMap)
		gm.Expect(err).To(gm.BeNil())

		sec := SecretArg{
			Name: "my-key",
			Type: "secret",
		}

		gm.Expect(named.Name).To(gm.Equal("test"))
		gm.Expect(named.Type).To(gm.Equal("tokenization"))
		gm.Expect(named.Args["key"]).To(gm.Equal(sec))
		gm.Expect(named.Args["nonSecret"]).To(gm.Equal(10))
	})
}
