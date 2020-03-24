// package config provides structs for holding, parsing, writing, and working
// with local configuration data for the command line tool.
package config

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"sigs.k8s.io/yaml"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// All files and folders that contain cape cli configuration must only be
// readable and writable by the owner. This is to ensure that another user on
// the system _cannot_ get access to the users `auth_token` for any of their
// clusters.
const requiredPermissions = 0600

var (
	ErrNoCluster   = errors.New(MissingConfigCause, "No cluster has been configured")
	ErrMissingHome = errors.New(InvalidEnvCause, "The $HOME environment variable could not be found")
)

// Default returns a Config struct with the default values set
func Default() *Config {
	return &Config{
		Version: 1,
		UI: UI{
			Colors:     true,
			Animations: true,
		},
		Context:  &Context{},
		Clusters: []Cluster{},
	}
}

// UI represents the configuration settings for how data is displayed
type UI struct {
	Colors     bool `json:"colors"`
	Animations bool `json:"animations"`
}

// Config represents the configuration settings for the command line
type Config struct {
	Version  int       `json:"version"`
	Context  *Context  `json:"context,omitempty"`
	Clusters []Cluster `json:"clusters,omitempty"`
	UI       UI        `json:"ui"`
}

// Write writes the configuration file out to the globally configured and
// derived location.
func (c *Config) Write() error {
	if err := c.Validate(); err != nil {
		return err
	}

	folderPath, err := FolderPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(folderPath, requiredPermissions)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	cfgPath, err := Path()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgPath, b, requiredPermissions)
}

// Print writes the configuration out the given stream
func (c *Config) Print(w io.Writer) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return nil
	}

	return nil
}

// Cluster returns the current cluster, if no cluster is set, cluster will be n
func (c *Config) Cluster() (*Cluster, error) {
	if c.Context.Cluster.String() == "" {
		return nil, ErrNoCluster
	}

	for _, cluster := range c.Clusters {
		if cluster.Label == c.Context.Cluster {
			return &cluster, nil
		}
	}

	return nil, errors.New(InvalidConfigCause, "The key 'context.cluster' is set but the cluster does not exist")
}

// Validate returns an error if the config is invalid
func (c *Config) Validate() error {
	if c.Version != 1 {
		return errors.New(InvalidVersionCause, "Expected a config version of 1")
	}

	if err := c.Context.Validate(); err != nil {
		return err
	}

	for _, cluster := range c.Clusters {
		if err := cluster.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Context represents the context section of the command line config
type Context struct {
	Cluster primitives.Label `json:"cluster,omitempty"`
}

// Validate returns an error if the context is invalid
func (c *Context) Validate() error {
	return c.Cluster.Validate()
}

// Cluster represents configuration for a Cape cluster
type Cluster struct {
	AuthToken string           `json:"auth_token,omitempty"`
	URL       *url.URL         `json:"url"`
	Label     primitives.Label `json:"label"`
}

// Validate returns an error if the cluster configuration is invalid
func (c *Cluster) Validate() error {
	return c.Label.Validate()
}

// Path returns the path to local configuration yaml file.
func Path() (string, error) {
	base, err := FolderPath()
	if err != nil {
		return "", err
	}

	return path.Join(base, "config.yaml"), nil
}

// FolderPath returns the path to the local folder that holds user-space wide
// cape configuration
//
// TODO: Add support for XDG_CONFIG standard which can be found at
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func FolderPath() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", ErrMissingHome
	}

	return path.Join(home, ".cape"), nil
}

// Parse reads the given file path and returns a Config object or returns an
// error as to why the config could not have been read
func Parse() (*Config, error) {
	filePath, err := Path()
	if err != nil {
		return nil, err
	}

	src, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	cfg := Default()
	if os.IsNotExist(err) {
		return cfg, nil
	}

	if src.Mode().Perm() != requiredPermissions {
		return nil, errors.New(InvalidPermissionsCause, "Invalid permissions for file %s, must be %d", filePath, requiredPermissions)
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
