package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/database/types"
)

// TestEntity represents an internal Entity used exclusively for testing
type TestEntity struct {
	*Primitive
	Data string `json:"data"`
}

type TestInnerEntity struct {
	Field1 string `json:"field_1"`
	Field2 string `json:"field_2"`
	Field3 string `json:"field_3"`
}

type TestNestedEntity struct {
	*Primitive
	Inner *TestInnerEntity
}

// GetType returns the type of this entity
func (t *TestNestedEntity) GetType() types.Type {
	return types.TestNested
}

// GetType returns the type of this entity
func (t *TestEntity) GetType() types.Type {
	return types.Test
}

// NewTestEntity returns a new TestEntity struct
func NewTestNestedEntity(inner *TestInnerEntity) (*TestNestedEntity, error) {
	p, err := NewPrimitive(types.Test)
	if err != nil {
		return nil, err
	}

	e := &TestNestedEntity{
		Primitive: p,
		Inner: inner,
	}

	// XXX: Static time for the purposes of testing
	e.CreatedAt = time.Unix(0, 0).UTC()
	e.UpdatedAt = time.Unix(0, 0).UTC()

	ID, err := DeriveID(e)
	if err != nil {
		return nil, err
	}

	e.ID = ID
	return e, nil
}

// NewTestEntity returns a new TestEntity struct
func NewTestEntity(data string) (*TestEntity, error) {
	p, err := NewPrimitive(types.Test)
	if err != nil {
		return nil, err
	}

	e := &TestEntity{
		Primitive: p,
		Data:      data,
	}

	// XXX: Static time for the purposes of testing
	e.CreatedAt = time.Unix(0, 0).UTC()
	e.UpdatedAt = time.Unix(0, 0).UTC()

	ID, err := DeriveID(e)
	if err != nil {
		return nil, err
	}

	e.ID = ID
	return e, nil
}

// TestMutableEntity represents an internal Entity used exclusively for testing
type TestMutableEntity struct {
	*Primitive
	Data string `json:"data"`
}

// GetType returns the type of this entity
func (t *TestMutableEntity) GetType() types.Type {
	return types.TestMutable
}

// NewTestMutableEntity returns a new TestMutableEntity
func NewTestMutableEntity(data string) (*TestMutableEntity, error) {
	p, err := NewPrimitive(types.TestMutable)
	if err != nil {
		return nil, err
	}

	return &TestMutableEntity{
		Primitive: p,
		Data:      data,
	}, nil
}

type TestEncryptionEntity struct {
	*Primitive
	Data string `json:"data"`
}

// GetType returns the type of this entity
func (t *TestEncryptionEntity) GetType() types.Type {
	return types.TestMutable
}

type testEncryptionEntity struct {
	*TestEncryptionEntity
	Data *base64.Value `json:"data"`
}

func NewTestEncryptionEntity(data string) (*TestEncryptionEntity, error) {
	p, err := NewPrimitive(types.TestMutable)
	if err != nil {
		return nil, err
	}

	e := &TestEncryptionEntity{
		Primitive: p,
		Data:      data,
	}

	return e, nil
}

func (t *TestEncryptionEntity) Encrypt(ctx context.Context, codec crypto.EncryptionCodec) ([]byte, error) {
	data, err := codec.Encrypt(ctx, base64.New([]byte(t.Data)))
	if err != nil {
		return nil, err
	}

	return json.Marshal(testEncryptionEntity{
		TestEncryptionEntity: t,
		Data:                 data,
	})
}

func (t *TestEncryptionEntity) Decrypt(ctx context.Context, codec crypto.EncryptionCodec, data []byte) error {
	in := &testEncryptionEntity{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	unencrypted, err := codec.Decrypt(ctx, in.Data)
	if err != nil {
		return err
	}

	t.Primitive = in.Primitive

	t.Data = string(*unencrypted)
	return nil
}

func (t *TestEncryptionEntity) GetEncryptable() bool {
	return true
}
