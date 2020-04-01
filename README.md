# Cape

TODO

## Components

### CLI

TODO

### Controller

TODO

### Data Connector

TODO

## Getting Started

TODO

## Development

Cape requires the following software to build and run locally:

- **Go 1.14** You can install Go using the official instructions found [here](https://golang.org/doc/install) or use [gvm](https://github.com/moovweb/gvm) to manage multiple go installations.
- **Protobuf** You can install the Protobuf compiler (`protoc`) using instructions found [here](https://developers.google.com/protocol-buffers/docs/downloads). We suggest using `brew` or your local system package manager (e.g. `apt`, `yum`, etc) to manage your version of `protoc`.
- **Docker** You can install Docker by following the instructions found [here](https://docs.docker.com/install/). 

Once you've installed the pre-requisite software you should be able to install additional tooling using `make bootstrap` in the root of your checkout of Cape.

We're using [Github Actions](https://github.com/features/actions) to automate our CI/CD. To test basic CI locally you can run the following command:

```
make ci
```

If this passes than this should also pass on Github Actions.

To test Docker image building you can run:

```
make docker
```

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
