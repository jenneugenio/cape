package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/privacyai/primitives/types"
)

func TestDerive(t *testing.T) {
	gm.RegisterTestingT(t)

	e, err := NewTestEntity("yo")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(e.ID.Type()).To(gm.Equal(types.Test))
	gm.Expect(e.ID.Version()).To(gm.Equal(byte(idVersion)))

	gm.Expect(e.ID.String()).To(gm.Equal("3m0gaa32a13ee54mk5xyjhd4w8"))
}

func TestGenerate(t *testing.T) {
	gm.RegisterTestingT(t)

	e, err := NewTestMutableEntity("ha")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(e.ID.Version()).To(gm.Equal(uint8(1)))

	v, err := e.ID.Type()
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(v.Mutable()).To(gm.BeTrue())
}

func TestGenerateFail(t *testing.T) {
	gm.RegisterTestingT(t)

	e, err := NewTestEntity("yo")
	gm.Expect(err).To(gm.BeNil())

	ID, err := GenerateID(e)
	gm.Expect(err).To(gm.Equal(ErrNotMutable))
	gm.Expect(ID).To(gm.Equal(EmptyID))
}

func TestDecodeString(t *testing.T) {
	gm.RegisterTestingT(t)

	e, err := NewTestEntity("yo")
	gm.Expect(err).To(gm.BeNil())

	ID, err := DeriveID(e)
	gm.Expect(err).To(gm.BeNil())

	IDTwo, err := DecodeFromString(ID.String())
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(ID).To(gm.Equal(IDTwo))
}
