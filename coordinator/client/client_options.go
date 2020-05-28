package client

import (
	"encoding/json"
	"fmt"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/mitchellh/mapstructure"
	"io"
)

type SchemaOptions struct {
	BlobPath string
}

func (s *SchemaOptions) UnmarshalGQL(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(t, s); err != nil {
			return err
		}

		return nil
	default:
		return errors.New(SerializationCause, "Could not serialize schema options")
	}
}

func (s SchemaOptions) MarshalGQL(w io.Writer) {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(j))
}

type SourceOptions struct {
	WithSchema    bool
	SchemaOptions *SchemaOptions
}
