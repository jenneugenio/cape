# CLI

This module contains the code for adding commands to our CLI for the coordinator, connector and client CLI.

## Adding a new command

The first step to adding a command is familiarizing yourself with the [CLI Style Guide](style-guide.md). The style guide must be adhered to to the best of our abilities so that we have a great and consistent CLI experience.

When creating a new high level command, a new file should be created where the file name is the name of the command (e.g. `cape update` would result in a `cmd/update.go` file). All sub-commands can live in this file or be pulled out into other files as needed.

In each high-level file there will be `init` block that initializes the high-level command and any sub-commands. Once the commands are initialized the high-level command is append to the commands list which is then used by the `cmd/root.go` file to add the commands to the app.

```
func init() {
	startCmd := &cli.Command{
        ...
	}

	coordinatorCmd := &cli.Command{
        ...
		Subcommands: []*cli.Command{startCmd},
	}

	commands = append(commands, coordinatorCmd)
}
```

See `cmd/coordinator.go` to see a full example of our the above would look.

All flags should go in the `cmd/flags.go` file and all errors should go in the `cmd/errors.go` file.
