package worker

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/primitives"
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestWorker(t *testing.T) {
	g := gm.NewWithT(t)

	tokenStr := "2013q9ue1c0bfey22yjmm8dcn4,AeY6KfOe_IhFC0BHzP5EPHBodHRwOi8vbG9jYWxob3N0OjgwODA"
	err := os.Setenv("CAPE_TOKEN", tokenStr)

	token, err := auth.ParseAPIToken("2013q9ue1c0bfey22yjmm8dcn4,AeY6KfOe_IhFC0BHzP5EPHBodHRwOi8vbG9jYWxob3N0OjgwODA")
	g.Expect(err).To(gm.BeNil())

	dbURL, err := primitives.NewDBURL("postgres://postgres:dev@localhost:5432/cape?sslmode=disable")
	g.Expect(err).To(gm.BeNil())

	worker, err := NewWorker(&Config{
		Token: token,
		DatabaseURL: dbURL,
	})
	g.Expect(err).To(gm.BeNil())

	err = worker.Start()
	g.Expect(err).To(gm.BeNil())
}
