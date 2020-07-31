# CLI Style Guide

### About this Document

This guide contains rules for the look, feel, and functionality of our command
line interfaces. To the best of our ability, we should follow this guide at all
times.

### Table of Contents

- [CLI Style Guide](#cli-style-guide)
    - [About this Document](#about-this-document)
    - [Table of Contents](#table-of-contents)
    - [Configuration](#configuration)
      - [Location and Permissions](#location-and-permissions)
      - [Modifying and Listing](#modifying-and-listing)
    - [Accepting Input](#accepting-input)
      - [Arguments and Prompts](#arguments-and-prompts)
      - [Flags](#flags)
      - [Environment Variables](#environment-variables)
      - [Reading from stdin](#reading-from-stdin)
    - [Signals, Exit Codes and Erroring](#signals-exit-codes-and-erroring)
    - [Confirmations](#confirmations)
    - [Organizing Functionality](#organizing-functionality)
    - [Writing Help](#writing-help)
    - [Displaying Output](#displaying-output)
      - [stdout vs stderr](#stdout-vs-stderr)
      - [Printing Tables and Labels](#printing-tables-and-labels)
      - [Using Colors and Styles](#using-colors-and-styles)
    - [Standard Functionality](#standard-functionality)

### Configuration

Configuration is any data that must be persisted between command invocations
either to reduce repetition by the user or to enable certain behaviours.

An example of this configuration would be whether or not progress bars and
colors are enabled or data for connecting to clusters.

#### Location and Permissions

All configuration that a user would want to persistent between command
invocations should be stored in a standard location. This location should be
overridable using an environment variable.

If the folder does not exist then it should be automatically created with the
permissions set to only be readable/writable by the owning user (0600).

If any files within the folder contain sensitive information (a login token),
then the cli should exit with an error if the file is too permissive (no longer 0600).

#### Modifying and Listing

A user should be able to list and set any configuration through a command as
well as by modifying the configuration file directly.

If a user hasn't specified a value for a configuration variable then the
default should be displayed.

### Accepting Input

Input can be provided by a user in various different forms from an interactive
prompt all the way to a configuration file. It's important that we are
consistent in how we use these different sources so a user can reason whether
something would be a flag, positional argument, or something accepted through a
prompt.

**Arguments and Prompts vs. Flags, Environment Variables vs. Config**

A command line argument should be used for specifying the targeted entities to
run the command (e.g. specifying the name of a system user). If a required
input is not provided by the user then they should be prompted to supply the
value.

All other input should be supplied via a flag or an environment variable (e.g.
`--db-url` for specifying a database to connect to). A flag should only be required if
the command is invoked in a situation where stderr is not attached to a
terminal. Not all flags require a short form but must have a long form (e.g.
`--service-id`).

All common flags should have a corresponding environment variable, but not all
environment variables should have flags.

Configuration should only be used for values that are set and stored over many
different commands across many different invocations (e.g. whether or not to
enable progress bars).

**Context Aware - Prompting vs Erroring**

If a command is being executed directly by a user (stderr is attached to a
terminal) then they are prompted for input if it's not supplied. Otherwise, an
error is returned specifying that a required argument was not provided.

**Precedence**

The most specific method takes precedence in a situation where the same input
could be provided from different sources. This way, if a user has set an
environment variable or configuration option they could override it if they
desired.

1. Input provided via an interactive prompt - (most precedence)
2. Input provided via a command line argument
3. Input provided via a command line flag
4. Input provided via an environment variable
5. Input provided via configuration (as specified in the
   [configuration](#configuration) section) - (least precedence)

#### Arguments and Prompts

If there is a direct target of a command then the identifier for the
subject should be provided as a command line argument. Any other subsequent
information should be supplied through prompts.

If the user provides the target, it should be displayed in the form of a prompt
but the value should be autoaccepted (display as if they were prompted but
don't actually prompt).

If the user does not provide a required argument to a read command they
**could** be prompted to select the target from a list or to provide a name for
a new target they are creating.

When multiple inputs are required they should be accepted via command line
arguments except in cases where the input is a password. Passwords **must**
only ever be provided through a prompt, read in from a file specified using a
flag, or through an environment variable so they are not leaked into
`.bash_history`. A one time token is not considered a password, allowing it to
be specified as an argument, flag, or environment variable.

It's important to determine whether an input is a required value or modifying a
default value. Only required values should be provided via command line arguments
except for a read command which supplies search or selection functionality.

All of our prompts are built on-top of [promptui](https://github.com/manifoldco/promptui)
with all prompts being displayed on stderr instead of stdout.

If any signal is received when a prompt is active the program should exit.

**Example**: Deleting a user

Diana must be specified below because they are the target of the command.

```bash
$ cape users delete diana
```

**Example**: Create a user

When creating a user there is extra information required so the user is prompted for it.

```bash
$ cape users create bob
Name: bob
Enter Password: ********
Confirm Password: ********
```

#### Flags

Depending on whether a command is a read or write operation a flag will either
modifying the output of the operation or provide additional input when
modifying the state of the targeted entity (e.g. a user). For every flag
a corresponding environment variable is also accepted.

A flag may also be used to circumvent certain types of prompts such as
confirming whether or not to proceed with an action (e.g. do you really want to
delete this user?).

All flags must have a long form (e.g. `--db-url`) and can accept input (e.g.
`--db-url URL`) or act as booleans (e.g. `--verbose`). Common flags (used by many
commands or used frequently with a command) should have short codes
(e.g. `-y, --yes`) for usability.

Flags should represent general concepts instead of having many different flags
for choosing a value. For example, if choosing a data source type you'd specify
it with `--source-type TYPE` instead of `--hadoop`, `--postgres`.

A flag's short and long form along with any options or possible default values
should be displayed in the help for the specific command.

Global flags that apply to every command should be kept to a minimum and only used
if absolutely necessary. Any global flag should have a corresponding environment
variable. All global flags should be documented including their environment variable.

**Example**: Flag accepting multiple values with a default

```
   --source-type TYPE, -s TYPE     Override data source type (hadoop, file, postgres) (default: file)
```

#### Environment Variables

All environment variables should be namespaced under `CAPE` (e.g. `CAPE_DB_URL`).

To the best of our ability, environment variable names should not be reused to
represent multiple concepts and shouldn't be overly generalized.
In some cases this will not make sense or be feasible
such as `-y, --yes` which should be represented as `CAPE_CONFIRM`.

For boolean values, an environment variable is represented as `true` if it has
any value other than 0 or unset (e.g. `CAPE_CONFIRM=1` is the same as passing
`-y, --yes`).

#### Reading from stdin

If a command supports reading the contents of a file in via stdin then it
should be documented in the help. There should also be a corresponding
optional argument allowing a user to specify the path to a file to import from.

### Signals, Exit Codes and Erroring

**Signals**

All signals received by Cape CLI commands should immediately result in the
graceful exit of the program. An error should be returned to the user
indicating why the program exited. The program should exit with a exit code of
1.

If a `SIGHUP` is received by a daemon process (e.g. coordinator)
then it should re-read any configuration files and restart to the best of it's ability.
Otherwise, a log line should be produced indicating that the `SIGHUP` was ignored.

**Exit Codes**

If the invoked command could not successfully complete for any reason the
program should terminate with an exit code of 1. Successful commands should
terminate with an exit code of 0.

**Errors**

All errors (including usage errors) must be printed to stderr. These messages
should be explicit, clear, and useful to the end-user. For example, if the user
provided a name of an entity that does not exist they should be told that no
entity exists for the given name. If the user does not provide a required
parameter or argument then the usage information for the command along with a
clear error should be provided.

To the best of our ability, error messages from the server should only be
displayed to the user if they're concise and actionable. Otherwise, a general
error message should be displayed (think equivalent to a 500 Error).

In the case of a running background agent, detailed error logs should be
written to an errors log for debugging purposes.

**Example**: Receiving a signal

```
$ cape coordinator start &
Cape coordinator is now listening at https://localhost:8080.
$ kill -HUP $!
Aborted: Received a 'SIGHUP'.
```

**Example**: Usage error

```
$ cape coordinator
NAME:
   cape coordinator - Control access to your data with the Cape Coordinator

USAGE:
   cape coordinator command [command options] [arguments...]

COMMANDS:
   start
   help, h  Shows a list of commands or help for one command
```

### Confirmations

Any command that performs a destructive activity or potentially dangerous
activity should contain a confirmation prompt. If a command is invoked in a
non-attached terminal then the callee must provide the `-y, --yes` flag.

If the user does not confirm the command via the prompt or through a flag then
an error should be returned.

### Organizing Functionality

To the best of our ability, we try to organize commands such that they are
`cape <subject> <action> <target>` where the subject is an entity such as
`users`, an `action` is the change the user is requesting (e.g. `create` a new user),
and `target` is the instance of the entity being created or modified (e.g. `diana` is
the user to be created).

Commands that have subcommands should not be callable, instead they are ways to
group functionality together. This way there are no surprising side effects of
running `cape coordinator`.

In some cases we may want to pull commonly used commands to the top level for
ease of discoverability. For instance, the ability to login or logout of your
account via `cape login` and `cape logout` or to list help via `cape help <command>`.
This is left to the discretion of the implementer, however, we
want to try to keep the top level help as concise as possible with the
understanding that `cape help` will be a first touch point for many users.

### Writing Help

The help command is the first place any user will go when they're figuring out
how to accomplish a task on the CLI. As such, it's super important that we
write readable and concise documentation. To that end, we
have established a few rules when it comes to writing documentation.

- Every command should have a clear and presice sentence that describes exactly
  *why* you'd use the command.
- Every command should have usage information that displays the name with the
  required and optional arguments. The usage should also describe where flags
  should be provided (e.g. after the command but before arguments). Required
  arguments should be documented using `<name>` while optional should be
  documented using `[-options]`. For example: `cape coordinator [-options] <name>`.
- A lengthier description should also be displayed when a user is inquirying
  about a specific command (not a command grouping, but an actual command that
  results in an action) that details what the command does and why a user would
  use it.
- If possible a link to any more detailed documentation related to the command.
- Any flags that are available for use when invoking a command should be
  displayed when listing help information (including global flags). Included in
  the help information should be the short and long form of the flag, whether
  it takes any input, a description explaining what the flag does, the default
  value, and the environment variable if relevant.

### Displaying Output

We strive to have our CLI output as easy as possible to consume, a user should
never have to struggle to find the information they're seeking. To achieve this
we want our output to be as consitent as possible across commands, clear and
concise as possible, and leverage the limited UI elements available (e.g.
colouring, tables, etc.) to draw attention to the most important information.

However, our UI styling should be smart and contextually relevant to make
viewing and scripting as simple as possible. For example, if the output is
being forwarded to a file (e.g. not attached to a terminal) it shouldn't be
colored.

#### stdout vs stderr

Only the successful result of a command should be printed to stdout. All
progress messages, prompts, and error output should be printed to stderr.  This
way a user can pipe output (such as a table) directly to a file, if an error
occurs, then the output file would be empty and the program would terminate
with a non-zero exit code. They'd still be prompted for any input or see any
progress/error messages as long as stderr was attached to a terminal.

In the case that stderr is not attached to a terminal, all progress and error
messages would still be printed but prompts would be automatically filled in
according with our guidelines around [Accepting Input](#accepting-input).

#### Printing Tables and Labels

When printing output or tables use tabs instead of spaces for alignment. This
makes it easy for anyone to write a script to parse any output and enables them
to control the tab size in their terminal. For a table, both the header and
data itself should be written to stdout.

In Golang, you can use an [ansiterm](https://github.com/juju/ansiterm) tab
writer for managing formatting, spacing, and other complexities to ease the
maintenance burden.

**Example**: Listing Users

```bash
$ cape users list
Name              Created            Roles
bob               March 9th, 2020    data-scientist, admin
```

#### Using Colors and Styles

Using ANSI escape codes for colouring and styling gives us tools to make
consuming information, differentiationg between elements, and drawing attention
to specific information easier. By coloring certain elements the same way
across the board we can reduce the cognitive overhead to find information, a
user can rely on their intuitive thinking skills once they've learned the
mental model.

All colors and styles should be disabled when the output is not being displayed
to a terminal. A user should also be able to disable colours across the board
through the CLI configuration.

All table headings and labels (e.g. `Name`) should be bolded. Any form of
status or state information should be color coded (red, green, gray) to
represent whether the entity is available, in a transition state, or not
available.

The use of other forms of styling and coloring is left to the discretion of the
implementer. However, it's important to keep in mind that we want to use color
and styling sparsely - we don't want the user to be overwhelmed or to undermine
the user of styling for more important pieces of information.

### Standard Functionality

All CLI tools should have the following functionality:

- The ability to list the documentation for a command using either `cape help
  <command> <?sub-command>`, `cape <command> <?sub-command> -h`, or
  `vape <command> <?sub-command> --help`
- The ability to obtain the version information via `cape version`, `cape -v`,
  or `cape --version` which returns the current version and the date it was
  built.
- All CLI tools should be distributed via popular package managers for each
  operating system (e.g. brew, apt, yum, etc) and zipped binaries for each
  major operating system and architecture should be available as well.

