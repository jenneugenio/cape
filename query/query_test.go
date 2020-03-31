package query

import (
	gm "github.com/onsi/gomega"
	"testing"
)

func TestQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("It cannot perform star commands", func(t *testing.T) {
		//q, err := Parse("SELECT * from transactions")
		//gm.Expect(err).To(gm.BeNil())
	})

	t.Run("It redacts a field you cannot access", func(t *testing.T) {
		gm.RegisterTestingT(t)

		r, err := NewRule("transactions", Deny, "processor")
		gm.Expect(err).To(gm.BeNil())
		p, err := NewPolicy(r)
		gm.Expect(err).To(gm.BeNil())

		q, err := Parse("SELECT processor, card_number, value FROM transactions")
		gm.Expect(err).To(gm.BeNil())

		q, err = q.Rewrite(p)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(q.Raw()).To(gm.Equal("SELECT card_number, value FROM transactions"))
	})

}
