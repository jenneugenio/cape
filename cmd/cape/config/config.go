// package config provides structs for holding, parsing, writing, and working
// with local configuration data for the command line tool.
package config

import (
	"fmt"
	"github.com/capeprivacy/cape/models"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/manifoldco/go-base64"
	"sigs.k8s.io/yaml"

	"github.com/capeprivacy/cape/coordinator"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// All files and folders that contain cape cli configuration must only be
// readable and writable by the owner. This is to ensure that another user on
// the system _cannot_ get access to the users `auth_token` for any of their
// clusters.
const requiredPermissions = 0700

var (
	ErrNoCluster   = errors.New(MissingConfigCause, "Missing cluster configuration; please set via 'cape config clusters use'")
	ErrMissingHome = errors.New(InvalidEnvCause, "The $HOME environment variable could not be found")
	ErrUserInfo    = errors.New(InvalidEnvCause, "Unable to retrieve info about current user")
)

// Default returns a Config struct with the default values set
func Default() *Config {
	return &Config{
		Version: 1,
		UI: &UI{
			Colors:     true,
			Animations: true,
		},
		Context:  &Context{},
		Clusters: []*Cluster{},
	}
}

// UI represents the configuration settings for how data is displayed
type UI struct {
	Colors     bool `json:"colors"`
	Animations bool `json:"animations"`
}

// Config represents the configuration settings for the command line
type Config struct {
	Version  int        `json:"version"`
	Context  *Context   `json:"context,omitempty"`
	Clusters []*Cluster `json:"clusters,omitempty"`
	UI       *UI        `json:"ui"`
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

// AddCluster attempts to add the cluster to the configuration
func (c *Config) AddCluster(label models.Label, url *models.URL, authToken string) (*Cluster, error) {
	for _, cluster := range c.Clusters {
		if cluster.Label == label {
			return nil, errors.New(ExistingClusterCause, "A cluster labeled %s already exists", label)
		}
	}

	cluster := &Cluster{
		Label:     label,
		URL:       url,
		AuthToken: authToken,
	}

	if err := cluster.Validate(); err != nil {
		return nil, err
	}

	c.Clusters = append(c.Clusters, cluster)
	return cluster, nil
}

// RemoveCluster attempts to remove the cluster from configuration
func (c *Config) RemoveCluster(label models.Label) error {
	idx, cluster := c.findCluster(label)
	if cluster == nil {
		return errors.New(ClusterNotFoundCause, "No cluster named '%s' exists", label)
	}

	c.Clusters = append(c.Clusters[:idx], c.Clusters[idx+1:]...)
	return nil
}

// Cluster returns the current cluster, if no cluster is set, cluster will be n
func (c *Config) Cluster() (*Cluster, error) {
	if c.Context.Cluster.String() == "" {
		return nil, ErrNoCluster
	}

	_, cluster := c.findCluster(c.Context.Cluster)
	if cluster != nil {
		return cluster, nil
	}

	return nil, errors.New(InvalidConfigCause, "The key 'context.cluster' is set but the cluster does not exist")
}

// HasCluster returns true if the provided label exists, false otherwise
func (c *Config) HasCluster(clusterLabel models.Label) bool {
	idx, _ := c.findCluster(clusterLabel)
	return idx > -1
}

// GetCluster returns the cluster by its label
func (c *Config) GetCluster(clusterStr string) (*Cluster, error) {
	label := models.Label(clusterStr)

	_, cluster := c.findCluster(label)
	if cluster != nil {
		return cluster, nil
	}

	return nil, errors.New(InvalidConfigCause, "The key 'context.cluster' is set but the cluster does not exist")
}

func (c *Config) findCluster(label models.Label) (int, *Cluster) {
	for i, cluster := range c.Clusters {
		if cluster.Label == label {
			return i, cluster
		}
	}

	return -1, nil
}

// Use sets the given label as the current cluster or to unset the current cluster by passing an empty label
func (c *Config) Use(label models.Label) error {
	if label == "" {
		c.Context.Cluster = ""
		return nil
	}

	var target *Cluster
	for _, cluster := range c.Clusters {
		if cluster.Label == label {
			target = cluster
		}
	}

	if target == nil {
		return errors.New(ClusterNotFoundCause, "A cluster labeled '%s' does not exist", label)
	}

	c.Context.Cluster = target.Label
	return nil
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
	Cluster models.Label `json:"cluster,omitempty"`
}

// Validate returns an error if the context is invalid
func (c *Context) Validate() error {
	if c.Cluster == "" {
		return nil
	}
	return nil
}

// Cluster represents configuration for a Cape cluster
type Cluster struct {
	AuthToken string       `json:"auth_token,omitempty"`
	URL       *models.URL  `json:"url"`
	Label     models.Label `json:"label"`
	CertFile  string       `json:"tls_cert,omitempty"`
}

// Validate returns an error if the cluster configuration is invalid
func (c *Cluster) Validate() error {
	if c.URL == nil {
		return errors.New(InvalidConfigCause, "Missing url property for '%s' cluster", c.Label)
	}

	return c.URL.Validate()
}

// GetURL parses the url and returns it
func (c *Cluster) GetURL() (*models.URL, error) {
	if c.URL == nil {
		return nil, errors.New(InvalidConfigCause, "Missing url property for '%s' cluster", c.Label)
	}

	return c.URL, nil
}

// Token parses the auth token and returns the base64 value
func (c *Cluster) Token() (*base64.Value, error) {
	if c.AuthToken == "" {
		return nil, nil
	}

	return base64.NewFromString(c.AuthToken)
}

// SetToken sets the token on the cluster
func (c *Cluster) SetToken(token *base64.Value) {
	if token == nil {
		c.AuthToken = ""
		return
	}

	c.AuthToken = token.String()
}

// String completes the Stringer interface
func (c *Cluster) String() string {
	return fmt.Sprintf("%s (%s)", c.Label, c.URL.String())
}

// Transport returns a configured coordinator transport for this cluster.
// This can be used with a client
func (c *Cluster) Transport() (coordinator.ClientTransport, error) {
	clusterURL, err := c.GetURL()
	if err != nil {
		return nil, err
	}

	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	return coordinator.NewHTTPTransport(clusterURL, token, c.CertFile), nil
}

// Path returns the path to local configuration yaml file.
func Path() (string, error) {
	base, err := FolderPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, "config.yaml"), nil
}

// FolderPath returns the path to the local folder that holds user-space wide
// cape configuration
//
// TODO: Add support for XDG_CONFIG standard which can be found at
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func FolderPath() (string, error) {
	capeHome := os.Getenv("CAPE_HOME")
	if capeHome != "" {
		return filepath.Clean(capeHome), nil
	}

	user, err := user.Current()
	if err != nil {
		return "", ErrUserInfo
	}

	return filepath.Join(user.HomeDir, ".cape"), nil
}

// Parse reads the given file path and returns a Config object or returns an
// error as to why the config could not have been read
func Parse() (*Config, error) {
	filePath, err := Path()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	cfg := Default()
	if os.IsNotExist(err) {
		return cfg, nil
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
