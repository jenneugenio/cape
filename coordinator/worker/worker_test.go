package worker

import (
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestWorker(t *testing.T) {
	g := gm.NewWithT(t)

	tokenStr := "2013q9ue1c0bfey22yjmm8dcn4,AeY6KfOe_IhFC0BHzP5EPHBodHRwOi8vbG9jYWxob3N0OjgwODA"

	err := os.Setenv("CAPE_TOKEN", tokenStr)
	g.Expect(err).To(gm.BeNil())

	worker, err := NewWorker()
	g.Expect(err).To(gm.BeNil())

	err = worker.Start()
	g.Expect(err).To(gm.BeNil())
}
