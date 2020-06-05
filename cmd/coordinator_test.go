package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"sigs.k8s.io/yaml"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator"
)

func TestCoordinatorConfiguration(t *testing.T) {
	gm.RegisterTestingT(t)

	tmpDir, err := ioutil.TempDir("", "cape-test-tempdir")
	gm.Expect(err).To(gm.BeNil())
	defer os.RemoveAll(tmpDir)

	file := path.Join(tmpDir, "config.yaml")

	url := "postgres://user:pass@host:5432/database"
	os.Setenv("CAPE_DB_URL", url)
	defer os.Unsetenv("CAPE_DB_URL")

	t.Run("Can generate a config file", func(t *testing.T) {
		app, _ := NewHarness(nil)

		err = app.Run([]string{
			"cape", "coordinator", "configure",
			"--out", file,
			"--port", "8080",
		})
		gm.Expect(err).To(gm.BeNil())

		cfg := &coordinator.Config{}
		by, err := ioutil.ReadFile(file)
		gm.Expect(err).To(gm.BeNil())

		err = yaml.Unmarshal(by, cfg)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(cfg.Port).To(gm.Equal(8080))
		gm.Expect(cfg.Version).To(gm.Equal(1))
		gm.Expect(len(*cfg.RootKey)).To(gm.Equal(32))
		gm.Expect(cfg.DB.Addr.String()).To(gm.Equal(url))
	})

	t.Run("Can generate a config map", func(t *testing.T) {
		var configMap map[string]interface{}

		app, _ := NewHarness(nil)

		err = app.Run([]string{
			"cape", "coordinator", "configure",
			"--out", file,
			"--port", "8080",
			"--format", "config-map",
		})
		gm.Expect(err).To(gm.BeNil())

		by, err := ioutil.ReadFile(file)
		gm.Expect(err).To(gm.BeNil())

		err = yaml.Unmarshal(by, &configMap)
		gm.Expect(err).To(gm.BeNil())

		yml := configMap["data"].(map[string]interface{})["coordinator-config.yaml"].(string)
		cfg := &coordinator.Config{}
		err = yaml.Unmarshal([]byte(yml), cfg)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(cfg.Port).To(gm.Equal(8080))
		gm.Expect(cfg.Version).To(gm.Equal(1))
		gm.Expect(len(*cfg.RootKey)).To(gm.Equal(32))
		gm.Expect(cfg.DB.Addr.String()).To(gm.Equal(url))
	})

	t.Run("Can generate a config map", func(t *testing.T) {
		var configMap map[string]interface{}

		app, _ := NewHarness(nil)

		err = app.Run([]string{
			"cape", "coordinator", "configure",
			"--out", file,
			"--port", "8080",
			"--format", "secret",
		})
		gm.Expect(err).To(gm.BeNil())

		by, err := ioutil.ReadFile(file)
		gm.Expect(err).To(gm.BeNil())

		err = yaml.Unmarshal(by, &configMap)
		gm.Expect(err).To(gm.BeNil())

		b64 := configMap["data"].(map[string]interface{})["coordinator-config.yaml"].(string)

		yml, err := base64.StdEncoding.DecodeString(b64)
		gm.Expect(err).To(gm.BeNil())

		cfg := &coordinator.Config{}
		err = yaml.Unmarshal(yml, cfg)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(cfg.Port).To(gm.Equal(8080))
		gm.Expect(cfg.Version).To(gm.Equal(1))
		gm.Expect(len(*cfg.RootKey)).To(gm.Equal(32))
		gm.Expect(cfg.DB.Addr.String()).To(gm.Equal(url))
	})
}
