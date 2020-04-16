package query

import (
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	label, err := primitives.NewLabel("my-data")
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can parse a valid query", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "SELECT * from transactions")
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Errors on an invalid query", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "jdksajdksajdkldklasj")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("You cannot pass a join", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "select * from transactions join people on transactions.person_id = people.id")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Only accepts selects", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "delete from transactions")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Tracks the collection", func(t *testing.T) {
		gm.RegisterTestingT(t)

		q, err := New(label, "SELECT * from transactions")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(q.Collection()).To(gm.Equal("transactions"))
	})
}
