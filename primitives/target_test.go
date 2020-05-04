package primitives

import (
	"fmt"
	gm "github.com/onsi/gomega"
	"testing"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestTarget(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Valid target", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := NewTarget("records:collection.transactions")
		gm.Expect(err).To(gm.BeNil())
	})

	invalids := []string{
		"hello",
		"wow:cool",
		"this.shouldnt.work",
		"invalidtype:hmm.okay",
	}

	for _, invalid := range invalids {
		t.Run(fmt.Sprintf("Invalid target: %s", invalid), func(t *testing.T) {
			gm.RegisterTestingT(t)
			_, err := NewTarget(invalid)
			gm.Expect(errors.CausedBy(err, InvalidTargetCause)).To(gm.BeTrue())
		})
	}
}

func TestTargetCollection(t *testing.T) {
	gm.RegisterTestingT(t)
	target, err := NewTarget("records:mycollection.transactions")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(target.Collection()).To(gm.Equal(Collection("mycollection")))
}

func TestTargetEntity(t *testing.T) {
	gm.RegisterTestingT(t)
	target, err := NewTarget("records:mycollection.transactions")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(target.Entity()).To(gm.Equal(Entity("transactions")))
}

func TestCollectionWildcard(t *testing.T) {
	gm.RegisterTestingT(t)
	target, err := NewTarget("records:*")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(target.Type()).To(gm.Equal(Records))
	gm.Expect(target.Collection().String()).To(gm.Equal("*"))
	gm.Expect(target.Entity().String()).To(gm.Equal("*"))
}

func TestEntityWildcard(t *testing.T) {
	gm.RegisterTestingT(t)
	target, err := NewTarget("records:mycollection.*")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(target.Type()).To(gm.Equal(Records))
	gm.Expect(target.Collection().String()).To(gm.Equal("mycollection"))
	gm.Expect(target.Entity().String()).To(gm.Equal("*"))
}

func TestInternalCollection(t *testing.T) {
	gm.RegisterTestingT(t)
	target, err := NewTarget("users:*")
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(target.Type().String()).To(gm.Equal("users"))
	gm.Expect(target.Collection().String()).To(gm.Equal("*"))
}
