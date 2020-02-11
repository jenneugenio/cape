package partyerrors

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestNew(t *testing.T) {
	t.Run("return an error", func(t *testing.T) {
		gm.RegisterTestingT(t)

		e := NewMulti(UnknownCause, []string{"one message", "two messages"})
		err := ToError(e)

		gm.Expect(err.Cause.Category).To(gm.Equal(InternalServerErrorCategory))
		gm.Expect(err.Messages).To(gm.Equal([]string{"one message", "two messages"}))
		gm.Expect(err.Code()).To(gm.Equal(int32(500)))
		gm.Expect(err.StatusCode()).To(gm.Equal(500))
	})
}
