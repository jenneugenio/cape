# PrivacyAI

TODO

## Components

### CLI

TODO

### Controller

TODO

### Data Connector

TODO

## Development

PrivacyAI required go1.13.X to build and run the executable. You can follow the official instructions to install it [here](https://golang.org/doc/install) or use [gvm] to manage multiple go installations.

We're using [Github Actions](https://github.com/features/actions) to automate our CI/CD. To test basic CI locally you can run the following command:

```
make ci
```

If this passes than this should also pass on Github Actions.

To test Docker image building you can run:

```
make docker
```