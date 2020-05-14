package main

import (
	"io/ioutil"
	"strconv"

	"github.com/manifoldco/go-base64"
	"github.com/urfave/cli/v2"
	"sigs.k8s.io/yaml"

	"github.com/kelseyhightower/envconfig"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/logging"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func init() {
	startCmd := &Command{
		Usage: "Start an instance of the Cape coordinator",
		Command: &cli.Command{
			Name:   "start",
			Action: startCoordinatorCmd,
			Flags: []cli.Flag{
				loggingTypeFlag(),
				loggingLevelFlag(),
				configFilesFlag(),
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

	key, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	cfg := &coordinator.Config{
		Version: 1,
		Port:    port,
		DB: &coordinator.DBConfig{
			Addr: dbURL,
		},
		InstanceID: "coordinator",
		RootKey:    base64.New(key[:]),
	}

	err = cfg.Write(out)
	if err != nil {
		return err
	}

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	err = u.Notify(ui.Warn, "%s contains sensitive information ensure to keep it safe and secret.", out)
	if err != nil {
		return err
	}

	return u.Template("Cape coordinator configuration generated. Run `cape coordinator start --file {{ . }}` to continue.", out)
}

func startCoordinatorCmd(c *cli.Context) error {
	cfg, err := getConfig(c)
	if err != nil {
		return err
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

	ctrl, err := coordinator.New(cfg, logger)
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
