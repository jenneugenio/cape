package transformations

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"golang.org/x/crypto/scrypt"
)

const (
	n           = 2
	r           = 8
	p           = 1
	keyLen      = 32
	maxTokenLen = 64
)

type TokenizationTransform struct {
	field   string
	key     []byte
	maxSize int
}

func (t *TokenizationTransform) tokenizeBytes(x []byte) (string, error) {
	// Set to 32 bits for sha256.
	// See doc for the other parameters: https://godoc.org/golang.org/x/crypto/scrypt
	hash, err := scrypt.Key(x, t.key, n, r, p, keyLen)
	hashHex := hex.EncodeToString(hash)

	size := len(hashHex)
	if size > t.maxSize {
		size = t.maxSize
	}
	return hashHex[0:size], err
}

func (t *TokenizationTransform) Transform(schema *proto.Schema, input *proto.Record) error {
	field, err := GetField(schema, input, t.field)
	if err != nil {
		return err
	}

	output := &proto.Field{}
	switch ty := field.GetValue().(type) {
	case *proto.Field_String_:
		res, err := t.tokenizeBytes([]byte(ty.String_))
		if err != nil {
			return err
		}
		output.Value = &proto.Field_String_{String_: res}
	case *proto.Field_Bytes:
		res, err := t.tokenizeBytes(ty.Bytes)
		if err != nil {
			return err
		}
		output.Value = &proto.Field_Bytes{Bytes: []byte(res)}
	}
	return SetField(schema, input, output, t.field)
}

func (t *TokenizationTransform) Initialize(args Args) error {
	key := make([]byte, keyLen)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}
	t.key = key

	maxSize, found, err := args.LookupFloat64("maxSize")
	if err != nil {
		return err
	}
	if found {
		if maxSize < 0 {
			return errors.New(UnsupportedType, "Unsupported max size: must be positive integer")
		}

		t.maxSize = int(maxSize)
	}

	return nil
}

func (t *TokenizationTransform) Validate(args Args) error {
	maxSize, found, err := args.LookupFloat64("maxSize")
	if err != nil {
		return err
	}
	if found && maxSize < 0 {
		return errors.New(UnsupportedType, "Unsupported max size: must be positive integer")
	}

	return nil
}

func (t *TokenizationTransform) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_BYTEA,
		proto.FieldType_CHAR,
		proto.FieldType_TEXT,
		proto.FieldType_VARCHAR,
	}
}

func (t *TokenizationTransform) Function() string {
	return "tokenization"
}

func (t *TokenizationTransform) Field() string {
	return t.field
}

func NewTokenizationTransform(field string) (Transformation, error) {
	t := &TokenizationTransform{
		field:   field,
		maxSize: maxTokenLen,
	}
	return t, nil
}
