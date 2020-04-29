package framework

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestVersionHandler(t *testing.T) {
	gm.RegisterTestingT(t)

	req := httptest.NewRequest("GET", "http://my.cape.com", nil)
	w := httptest.NewRecorder()

	handler := VersionHandler("cape-hi")
	handler.ServeHTTP(w, req)

	resp := w.Result()
	gm.Expect(resp.StatusCode).To(gm.Equal(http.StatusOK))
	gm.Expect(resp.Header.Get("Content-Type")).To(gm.Equal("application/json"))

	v := &VersionResponse{}
	err := json.NewDecoder(resp.Body).Decode(v)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(v.InstanceID).To(gm.Equal("cape-hi"))
	gm.Expect(v.Version).To(gm.Equal("0.0.0"))
	gm.Expect(v.BuildDate).To(gm.Equal("never"))
}
