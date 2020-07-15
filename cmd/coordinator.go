package main

import (
	stdbase64 "encoding/base64"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Masterminds/sprig"
	"github.com/urfave/cli/v2"
	"sigs.k8s.io/yaml"

	"github.com/kelseyhightower/envconfig"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/mailer"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/logging"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"

	_ "github.com/capeprivacy/cape"
)

var typeRegistry = map[FormatType]string{
	K8Secret:    K8Secret.String(),
	K8ConfigMap: K8ConfigMap.String(),
}

func init() {
	startCmd := &Command{
		Usage: "Start an instance of the Cape coordinator",
		Variables: []*EnvVar{
			{
				Name:        "CAPE_CORS_ENABLE",
				Required:    false,
				Description: "Serve CORS headers in HTTP responses",
			},
			{
				Name:        "CAPE_CORS_ALLOW_ORIGIN",
				Required:    false,
				Description: "Specify the value of the CORS Allow Origin Header",
			},
		},
		Command: &cli.Command{
			Name:   "start",
			Action: startCoordinatorCmd,
			Flags: []cli.Flag{
				loggingTypeFlag(),
				loggingLevelFlag(),
				configFilesFlag(),
				instanceIDFlag(),
			},
		},
	}

	configureCmd := &Command{
		Usage:     "Generates a Cape coordinator configuration file",
		Variables: []*EnvVar{capeDBURLNotRequired},
		Examples: []*Example{
			{
				Example:     "cape coordinator configure",
				Description: "Generates a configuration file by prompting for the port and database url.",
			},
			{
				Example:     "cape coordinator configure --out my-config.yaml",
				Description: "Generates a configuration file and outputs it to my-config.yaml",
			},
			{
				Example: "CAPE_DB_URL=postgres://user:pass@host:5432/database cape coordinator configure --port 8080",
				Description: "Generates a configuration file without prompting " +
					"with the database url postgres://user:pass@host:5432/database and port 8080.",
			},
		},
		Command: &cli.Command{
			Name:   "configure",
			Action: configureCoordinatorCmd,
			Flags: []cli.Flag{
				configFileOutFlag(),
				portFlag("port", 0),
				formatFlag(),
			},
		},
	}

	coordinatorCmd := &Command{
		Usage: "Commands for starting and managing Cape coordinators.",
		Command: &cli.Command{
			Name:        "coordinator",
			Subcommands: []*cli.Command{startCmd.Package(), configureCmd.Package()},
		},
	}

	commands = append(commands, coordinatorCmd.Package())
}

func getConfig(c *cli.Context) (*coordinator.Config, error) {
	configs := c.StringSlice("file")

	config := &coordinator.Config{}
	for _, configFile := range configs {
		by, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(by, config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func getPort(c *cli.Context) (int, error) {
	validatePort := func(input string) error {
		port, err := strconv.Atoi(input)
		if err != nil {
			return err
		}

		if port > 65535 || port < 1 {
			return errors.New(InvalidPortCause, "Port must be between 1-65335")
		}

		return nil
	}

	question := "Please enter the port to run the coordinator on"

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	portStr, err := u.Question(question, validatePort)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(portStr)
}

func queryDBURL(c *cli.Context) (*primitives.DBURL, error) {
	validateDBURL := func(input string) error {
		_, err := primitives.NewDBURL(input)
		if err != nil {
			return err
		}

		return nil
	}

	question := "Enter the url of your database (e.g. postgres://user:pass@host:5432/database)"

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	urlStr, err := u.Question(question, validateDBURL)
	if err != nil {
		return nil, err
	}

	return primitives.NewDBURL(urlStr)
}

func configureCoordinatorCmd(c *cli.Context) error {
	out := c.String("out")

	port := c.Int("port")
	if port == 0 {
		p, err := getPort(c)
		if err != nil {
			return err
		}

		port = p
	}

	var dbURL *primitives.DBURL
	i := EnvVariables(c.Context, capeDBURLNotRequired)
	if i == nil {
		u, err := queryDBURL(c)
		if err != nil {
			return err
		}

		dbURL = u
	} else {
		dbURL = i.(*primitives.DBURL)
	}

	cfg, err := coordinator.NewConfig(port, dbURL)
	if err != nil {
		return err
	}

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	format := c.String("format")
	if format != "" {
		err := handleK8sFormat(cfg, format, out)
		if err != nil {
			return err
		}

		err = u.Notify(ui.Warn, "%s contains sensitive information ensure to keep it safe and secret.", out)
		if err != nil {
			return err
		}

		return u.Template("Cape coordinator configuration generated. Run `kubectl apply -f {{ . }}` to continue.", out)
	}

	err = cfg.Write(out)
	if err != nil {
		return err
	}

	err = u.Notify(ui.Warn, "%s contains sensitive information ensure to keep it safe and secret.", out)
	if err != nil {
		return err
	}

	return u.Template("Cape coordinator configuration generated. Run `cape coordinator start --file {{ . }}` to continue.", out)
}

func handleK8sFormat(cfg *coordinator.Config, format string, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return errors.New(CreateFileCause, "Unable to create file %s", out)
	}
	defer f.Close()

	by, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	switch FormatType(format) {
	case K8ConfigMap:
		tpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(configMapTemplate)
		if err != nil {
			return err
		}

		err = tpl.Execute(f, string(by))
		if err != nil {
			return err
		}
	case K8Secret:
		tpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(secretTemplate)
		if err != nil {
			return err
		}

		encoded := stdbase64.StdEncoding.EncodeToString(by)
		err = tpl.Execute(f, encoded)
		if err != nil {
			return err
		}
	}

	return nil
}

func startCoordinatorCmd(c *cli.Context) error {
	cfg, err := getConfig(c)
	if err != nil {
		return err
	}

	// By default, we always configure Cape to use Argon2ID
	// _unless_ we're writing unit tests.
	cfg.CredentialProducerAlg = primitives.Argon2ID

	instanceID, err := getInstanceID(c, "coordinator")
	if err != nil {
		return err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = instanceID
	}

	err = envconfig.Process("cape", cfg)
	if err != nil {
		return err
	}

	// TODO: Consider having the "logger" be configured by the server?
	logger, err := logging.Logger(c.String("logger"), c.String("log-level"), cfg.InstanceID.String())
	if err != nil {
		return err
	}

	// TODO: Enable proper configuration of the mailer including support for
	// sending email via SMPT over TLS. All email must be sent over TLS.
	mailer := &mailer.TestMailer{}
	ctrl, err := coordinator.New(cfg, logger, mailer)
	if err != nil {
		return err
	}

	server, err := framework.NewServer(cfg, ctrl, logger)
	if err != nil {
		return err
	}

	watcher, err := setupSignalWatcher(server, logger)
	if err != nil {
		return err
	}

	err = watcher.Start()
	if err != nil {
		return err
	}
	defer watcher.Stop()

	return server.Start(c.Context)
}

type FormatType string

func (f FormatType) String() string {
	return string(f)
}

const (
	K8Secret    FormatType = "secret"
	K8ConfigMap FormatType = "config-map"
)

// FormatTypes returns a map of a type to string representation
func FormatTypes() map[FormatType]string {
	return typeRegistry
}
