package types

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var testingType Type = 0xD03
var badType Type = 0xD03
var otherType Type = 0xD04
var overflowType Type = 0xFFF1

func TestRegister(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("can register and fetch", func(t *testing.T) {
		gm.RegisterTestingT(t)

		Register(testingType, "testing_type", false)

		et := Get("testing_type")
		gm.Expect(et).To(gm.Equal(testingType))
		gm.Expect(et.Mutable()).To(gm.Equal(false))
	})

	t.Run("get fail", func(t *testing.T) {
		gm.RegisterTestingT(t)
		gm.Expect(func() {
			Get("bleh")
		}).Should(gm.Panic())
	})

	t.Run("cannot register same byte representation", func(t *testing.T) {
		gm.RegisterTestingT(t)

		gm.Expect(func() {
			Register(badType, "woo_hoo", false)
		}).Should(gm.Panic())
	})

	t.Run("cannot register same name", func(t *testing.T) {
		gm.RegisterTestingT(t)

		gm.Expect(func() {
			Register(otherType, testingType.String(), false)
		}).Should(gm.Panic())
	})

	t.Run("overflow test", func(t *testing.T) {
		gm.RegisterTestingT(t)

		gm.Expect(func() {
			Register(overflowType, testingType.String(), false)
		}).Should(gm.Panic())
	})
}

func TestDecode(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("type exists", func(t *testing.T) {
		gm.RegisterTestingT(t)

		et, err := Decode(uint16(0xD03))
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(et).To(gm.Equal(testingType))
	})

	t.Run("type does not exist", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := Decode(uint16(0xDDD))
		gm.Expect(errors.FromCause(err, UnknownTypeCause)).To(gm.BeTrue())
	})

	t.Run("byte decode", func(t *testing.T) {
		gm.RegisterTestingT(t)

		ty := Get("testing_type")
		gm.Expect(ty).To(gm.Equal(testingType))

		up := ty.Upper()
		gm.Expect(up).To(gm.Equal(byte(13)))
		lo := ty.Lower()
		gm.Expect(lo).To(gm.Equal(byte(3)))

		ty2, err := DecodeBytes(up, lo)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(ty2).To(gm.Equal(testingType))
	})

	t.Run("mutable type", func(t *testing.T) {
		gm.RegisterTestingT(t)

		v := Get("test_mutable")
		gm.Expect(v.Mutable()).To(gm.BeTrue())
	})
}
