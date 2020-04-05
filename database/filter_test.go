package database

import (
	"testing"

	gm "github.com/onsi/gomega"
)

type tester string

func (t tester) String() string {
	return string(t)
}

func TestIn(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("in can uniqify values", func(t *testing.T) {
		out := In{"a", "b", "a"}
		gm.Expect(out.Uniqify()).To(gm.Equal(In{"a", "b"}))
	})

	t.Run("uniqify panics on non-string or stringer", func(t *testing.T) {
		in := In{1, 2, 3}
		gm.Expect(func() {
			in.Uniqify()
		}).To(gm.Panic())
	})

	t.Run("values returns stringified values", func(t *testing.T) {
		in := In{tester("hi"), tester("b")}
		gm.Expect(in.Values()).To(gm.Equal([]interface{}{"hi", "b"}))
	})

	t.Run("values returns string values", func(t *testing.T) {
		in := In{"a", "b"}
		gm.Expect(in.Values()).To(gm.Equal([]interface{}{"a", "b"}))
	})

	t.Run("values panics on non-string or stringer", func(t *testing.T) {
		in := In{1, 2}
		gm.Expect(func() {
			in.Values()
		}).To(gm.Panic())
	})

	t.Run("can create in from array of entities", func(t *testing.T) {
		eA, err := NewTestEntity("hello")
		gm.Expect(err).To(gm.BeNil())

		eB, err := NewTestEntity("yeswhynot")
		gm.Expect(err).To(gm.BeNil())

		in := []*TestEntity{eA, eB}
		out := InFromEntities(in, func(e interface{}) interface{} {
			return e.(*TestEntity).Data
		})

		gm.Expect(out).To(gm.Equal(In{eA.Data, eB.Data}))
	})
}
