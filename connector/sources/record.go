package sources

import "github.com/dropoutlabs/cape/connector/proto"

// Record contains the values decoded by our decoder
type Record struct {
	values []interface{}
}

// NewRecord decodes the incoming bytes given the schema and returns a Record
func NewRecord(schema *proto.Schema, data [][]byte) (*Record, error) {
	values, err := Decode(schema, data)
	if err != nil {
		return nil, err
	}

	return &Record{
		values: values,
	}, nil
}

// Values returns the underlying values
func (r *Record) Values() []interface{} {
	return r.values
}
