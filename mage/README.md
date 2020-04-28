# Magefile for Cape

This package contains all of the files used to construct the Cape
[Magefile](https://magefile.org) experience.

### Architecture

**Package `mage`**

The root of the package (`mage`) contains the individual components that wrap
our different dependencies and workflows. These components that been designed
such that they don't contain particulars about our environment, instead,
they're generic re-usable components that could be used to create many
different workflows.

Many of these components satisfy the [`mage.Dependency`](./dependencies.go)
interface which specifies a methods that can be used for managing the lifecycle
of things like Go, Docker, Protoc, and others.

In addition to those that implement `Dependency` there are also helpers like
the [`Artifacts`](./artifacts.go) tracker which makes it easy to track the
artifacts that _could_ be created during execution of a target.

Our goal is to decouple the core functionality from the
specifics of our workflow making it easy and possible to write automated tests
for our most critical workflows.

**Package `mage/targets`**

All of the targets are specified inside the [`mage/targets`](./targets)
package. This package contains the workflow specific code and documents which
artifacts are created by which commands.
