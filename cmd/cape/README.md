# CLI

This module contains the code for adding commands to our CLI for the coordinator and client CLI.

## Adding a new command

The first step to adding a command is familiarizing yourself with the [CLI Style Guide](style-guide.md).
The style guide must be adhered to to the best of our abilities so that we have a great and consistent CLI experience.

When creating a new high level command, a new file should be created where the file name is the name of the command
(e.g. `cape update` would result in a `cmd/update.go` file).
All sub-commands can live in this file or be pulled out into other files as needed.

In each high-level file there will be `init` block that initializes the high-level command and any sub-commands.
Once the commands are initialized the high-level command is append to the commands list which is then used by the `cmd/root.go`
file to add the commands to the app.

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

All flags should go in the `cmd/flags.go` file, all errors should go in the `cmd/errors.go` file, and
arguments in `cmd/arguments.go`.

## Implementing a new command

When you are implementing a command, you will declare a function that will be passed a `*cli.Context` variable.
There are various middlewares and helpers that you can use to make writing a new command easy.

### Accessing Flags

You can easily pull your flags off of the `cli.Context`. E.g. if you specified a `port` flag, you would
access it with `c.Int("port")`.

See `cmd/flags.go` for examples of declaring flags.

### Accessing Arguments

Similarly, if you declared a URL argument named `MyURLArgument`, you could access it like so

```go
URLValue, ok := Arguments(c.Context, MyURLArgument).(primitive.URL)
```

see `cmd/arguments.go` for examples on declaring arguments.

### Communicating with the server

Through our middleware, you can easily get a reference to a client that can talk to the cape coordinator.

This can be done with

```go
provider := GetProvider(c.Context)
client, err := provider.Client(c.Context)
```

### Printing to the screen

One thing to note, is that your command should not ever directly print output. E.g. you should never write `fmt.Printf`.
Instead, you should use the `ui` object. This will handle how to print what you want, where to print it based on user settings,
and more.

Again, the `UI` is automatically provided to you via middleware. You can get access to it like so

```go
provider := GetProvider(c.Context)
u := provider.UI(c.Context)
```

`cmd/tokens.go` has a good examples of different ways you can use the `ui` class.
