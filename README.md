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

In addition, if you plan on running tests or building Cape outside of a docker
container, you'll need to install:

- [protoc](https://developers.google.com/protocol-buffers/docs/downloads) (version 3.11.4+)

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
helm charts to get the coordinator, the connector and the database running. This can be done with the
following `mage` command:

```
$ mage local:deploy
```

You should see something like:

```
NAME: coordinator
LAST DEPLOYED: Fri May 15 16:23:30 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NAME: connector
LAST DEPLOYED: Fri May 15 16:23:36 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

to know it has completed successfully.

Once that command is done running there are a few more commands to get everything set up.

First, run setup to create an admin account, create the default roles and the default policies.

```
$ cape setup local http://localhost:8080
```

Next, create a service which can be used as the data connector. The connector won't actually be
running until you complete this step has it requires the kubernetes secret inserted in the last command.

```
$ cape services create --type data-connector --endpoint https://localhost:8081 service:dc@my-cape.com
$ export CONNECTOR_TOKEN=<TOKEN PRINTED OUT FROM THE LAST COMMAND>
$ kubectl create secret generic connector-secret --from-literal=token=$CONNECTOR_TOKEN
```

Get a token that you can provide the worker to be able to auth with the connector & coordinator
```
$ cape services create --type worker service:worker@my-cape.com
$ export WORKER_TOKEN=<TOKEN PRINTED OUT FROM THE LAST COMMAND>
$ kubectl create secret generic worker-secret --from-literal=token=$WORKER_TOKEN
```

Create a data source and point it to some test data located in the local cluster. This needs to be
explicitly linked with the data connector created above.

```
$ cape sources add --link service:dc@my-cape.com transactions postgres://postgres:dev@postgres-customer-postgresql:5434/customer
```

Create a policy so that the data is actually accessible. By default, access to data in the Cape is
blocked unless there is a policy saying you can access it.

```
$ cape policies attach --from-file examples/allow-specific-fields.yaml allow-specific-fields global
```

Try to pull the data! If this command succeeds you should be good to continue experimenting and testing
the system.

```
$ cape pull transactions "SELECT * FROM transactions"
```

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

