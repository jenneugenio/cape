package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/privacyai/primitives/types"
)

func TestDerive(t *testing.T) {
	gm.RegisterTestingT(t)

	e := &TestEntity{Data: "yo"}
	ID, err := DeriveID(e)
	gm.Expect(err).To(gm.BeNil())

	IDTwo, err := DeriveID(e)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(ID).To(gm.Equal(IDTwo))

	gm.Expect(ID.Type()).To(gm.Equal(types.Test))
	gm.Expect(ID.Version()).To(gm.Equal(byte(idVersion)))

	gm.Expect(ID.String()).To(gm.Equal("3m0hku24b4t4u03fv0bk36ec5w"))
}

func TestGenerate(t *testing.T) {
	gm.RegisterTestingT(t)

	e := &TestingMutableEntity{Data: "ha"}
	ID, err := GenerateID(e)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(ID.Version()).To(gm.Equal(uint8(1)))

	v, err := ID.Type()
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(v.Mutable()).To(gm.BeTrue())
}

func TestGenerateFail(t *testing.T) {
	gm.RegisterTestingT(t)

	e := &TestEntity{Data: "yo"}
	ID, err := GenerateID(e)
	gm.Expect(err).To(gm.Equal(ErrNotMutable))
	gm.Expect(ID).To(gm.Equal(EmptyID))
}

func TestDecodeString(t *testing.T) {
	gm.RegisterTestingT(t)

	e := &TestEntity{Data: "ha"}
	ID, err := DeriveID(e)
	gm.Expect(err).To(gm.BeNil())

	IDTwo, err := DecodeFromString(ID.String())
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(ID).To(gm.Equal(IDTwo))
}
