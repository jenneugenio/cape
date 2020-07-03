package scratch

import (
	gm "github.com/onsi/gomega"
	"testing"
)

func TestDoIt(t *testing.T) {
	gm.RegisterTestingT(t)

	DoIt()
}
