package worker

import (
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestWorker(t *testing.T) {
	g := gm.NewWithT(t)
	//
	//s, err := harness.NewStack(ctx)
	//g.Expect(err).To(gm.BeNil())
	//defer s.Teardown(ctx)
	//
	//err = s.Manager.CreateSource(ctx, s.ConnHarness.SourceCredentials(), s.Manager.Connector.ID)
	//g.Expect(err).To(gm.BeNil())
	//
	//token, err := s.CoordClient.CreateToken(ctx, nil)
	//g.Expect(err).To(gm.BeNil())

	//tokenStr, err := token.Marshal()
	//g.Expect(err).To(gm.BeNil())

	tokenStr := "2017a7nwa6xkkb8gr3kehxae8w,AXUwRmxZvoCVEWB0CdkGgylodHRwOi8vbG9jYWxob3N0OjgwODA"

	err := os.Setenv("CAPE_TOKEN", tokenStr)
	g.Expect(err).To(gm.BeNil())

	worker, err := NewWorker()
	g.Expect(err).To(gm.BeNil())

	err = worker.Start()
	g.Expect(err).To(gm.BeNil())
}
