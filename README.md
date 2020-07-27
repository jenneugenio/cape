# Cape

## Development

In the following section we describe how to setup your environment so you can
develop Cape itself locally including how to test, build, and deploy Cape to a
local kubernetes cluster.

### Getting Started

To get up and running with contributing to Cape or to run Cape locally using
our development tooling you will need to have the following dependencies
installed:

- [git](https://git-scm.com/) (version 2.0+)
- [docker](https://docs.docker.com/get-docker/) (version 18.0+)
- [golang](https://golang.org/doc/install) (version 1.14+)

Once you have the base dependencies available on your system you can bootstrap
your local environment by running:

```
$ go run bootstrap.go bootstrap
```

This will install any additional local dependencies and check that your system
has everything installed to build and deploy Cape to your local environment.

To see a list of all local development options please run `mage -l` in the root
of the repository.

At any time, you can run `mage check` to check the status of your local
development environment.

### Building

Once your development environment is bootstrapped you can build the Cape binary
by running the following command:

```
$ mage build:binary
```

This will output the cape binary to `bin/cape`.

### Testing

You can run the full test suite locally through mage by running `mage test:ci`.
For these tests to pass you must have an instance of postgres running that is
connectable and configured. See the list of testing commands via `mage -l`.

We're using [GitHub Actions](https://github.com/features/actions) to automate
our continus integration and delivery suite which invokes `mage test:ci` to
determine the state of the build.

### Local Deployment

We're using [kind](https://kind.sigs.k8s.io/) for deploying Cape into kubernetes locally.

To get started you need to first startup kind, build the docker images and deploy the
helm charts to get the coordinator and the database running. This can be done with the
following `mage` command:

```
$ mage local:setup
```

You should see something like:

```
NAME: coordinator
LAST DEPLOYED: Fri May 15 16:23:30 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

to know it has completed successfully.

This creates an admin account for you which has the email *cape_user@mycape.com*
and password *capecape*. If required, you can login again with:

```bash
$ CAPE_PASSWORD=capecape cape login --email cape_user@mycape.com
```

### Creating API Token

Logging in from cape-python requires an API token. An APIT token also optionally be used to log in via the CLI.
This token can be acquired by running:

```bash
$ cape tokens create
```

You should see an output like:

```
A token for cape_user@mycape.com has been created!

Token:    2015te16gduwfrq5reedt7ygjr,AR0GdmFMKzC42u7pA5sfYMgiQEGXp2aa6A

‼ Remember: Please keep the token safe and share it only over secure channels.
```

`2015te16gduwfrq5reedt7ygjr,AR0GdmFMKzC42u7pA5sfYMgiQEGXp2aa6A` can be copied and pasted to
where its needed.

### Create an Initial Project

Projects are used to manage which version of policy is applied to a dataset. To create a
first project run:

```bash
$ cape projects create first-project "Hello Project World"
```

Once the project is created we can add policy to it from a policy specification:

```bash
$ cape projects update --from-spec examples/perturb_ones_field.yaml first-project
```

This project can now be used to perturb actual data from cape-python.

## Troubleshooting

Here are some steps to try if you've installed `cape` and you're
having one of these problems:

### Fedora

#### Docker does not support cgroups v2

Fedora 31 migrated from cgroups v1 to v2, but Docker and Kubernetes
don't support cgroups v2. Update the kernel to use cgroups v1.

```sh
sudo dnf install -y grubby && \
  sudo grubby \
  --update-kernel=ALL \
  --args="systemd.unified_cgroup_hierarchy=0"
```

#### Unable to resolve host from inside the container

This is an
[issue](https://github.com/kubernetes-sigs/kind/issues/1547) that
surfaced in Fedora 32 when using kind.

Run `nmcli connection show` in your terminal to get the ethernet
interface device ID.

```sh
firewall-cmd --permanent --zone=trusted --add-interface=docker0
firewall-cmd --get-zone-of-interface=<your eth interface>
firewall-cmd --zone=<zone from above> --add-masquerade --permanent
firewall-cmd --reload
```

Destroy your local environment with `mage local:destroy` and run
through the setup process once more.

