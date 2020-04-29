package database

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestDerive(t *testing.T) {
	gm.RegisterTestingT(t)

	e, err := NewTestEntity("yo")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(e.ID.Type()).To(gm.Equal(types.Test))
	gm.Expect(e.ID.Version()).To(gm.Equal(byte(idVersion)))

	gm.Expect(e.ID.String()).To(gm.Equal("3m0kg05yh1wzag8qux42umjz84"))
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

	ID, err := GenerateID(types.Test)
	gm.Expect(errors.FromCause(err, NotMutableCause)).To(gm.BeTrue())
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

func TestGQL(t *testing.T) {
	t.Run("Can gql unmarshal", func(t *testing.T) {
		e, err := NewTestMutableEntity("ha")
		gm.Expect(err).To(gm.BeNil())

		var id ID

		err = id.UnmarshalGQL(e.ID.String())
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(id.String()).To(gm.Equal(e.ID.String()))
	})

	t.Run("GQL unmarshal throws validate error", func(t *testing.T) {
		id := EmptyID

		var newID ID

		err := newID.UnmarshalGQL(id.String())
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidIDCause))
	})
}