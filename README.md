# Cape

TODO

## Components

### CLI

TODO

### Coordinator

TODO

### Data Connector

TODO

## Getting Started

TODO

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

### Tilt

[Tilt](tilt.dev) can be used to test everything locally including Kubernetes deployment and application features. Getting started is fairly easy. On MacOS, we recommend using Docker for Desktop which comes with a Kubernetes installion. You must enable it though which can be done by following [this](https://docs.docker.com/docker-for-mac/#kubernetes) documentation. If Kubernetes is having trouble starting you may need to reset factory defaults for Docker for Desktop.

Install tilt with:

```
curl -fsSL https://raw.githubusercontent.com/windmilleng/tilt/master/scripts/install.sh | bash
```

Which can then be run with:

```
tilt up
```

By default, `tilt` starts a webserver which can be reached from your browser. Here you can manage tilt and watch the logs of its actions. There is a similar thing launched on CLI so if you want to launch without browser you can run:

```
tilt up --no-browser
```

`tilt` detects changes on the source code and automatically rebuilds the docker containers and relaunches the services and deployments so once you've run `tilt up` you shouldn't have to bring it down until you're done testing.

**Resources**

- https://docs.tilt.dev/welcome_to_tilt.html
- https://docs.tilt.dev/install.html
- https://docs.tilt.dev/onboarding_checklist.html

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

