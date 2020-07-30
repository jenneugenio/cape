# Contribution guide

Contributions are more than welcome and we're always looking for use cases and feature ideas!

This document helps you get started on:

- [Set up local development](#local-development)
- [Submitting a pull request](#submitting-a-pull-request)
- [Writing documentation](#writing-documentation)
- [Useful tricks](#useful-tricks)
- [Reporting a bug](#reporting-a-bug)
- [Asking for help](#asking-for-help)

## Local development

This section describes how to set up your local environment so you can develop Cape, including how to test, build, and deploy Cape to a local Kubernetes cluster.

### Prerequisites

- [git](https://git-scm.com/) 2.0+
- [Docker](https://docs.Docker.com/get-Docker/) 18.0+
- [Go](https://golang.org/doc/install) 1.14+
- [PostgreSQL](https://www.postgresql.org/) 11.0+ (to run the CI tests)
- Linux (any distro) or MacOS (10.15 - Catalina or greater)

### PostgreSQL Install

```shell
# Install PostgreSQL
# Refer to the guidance for your distribution: https://www.postgresql.org/download/

# Create a Cape user
createuser --createdb cape

# Create a Cape database
createdb -U cape cape
```

### Install and set up

Clone the repo:

```
git clone https://github.com/capeprivacy/cape.git
```

Bootstrap your local environment:

```
$ go run bootstrap.go bootstrap
```

This will install any additional local dependencies and check that your system
has everything installed to build and deploy Cape to your local environment.

To see a list of all local development options, run `mage -l` in the root
of the repository.

At any time, you can run `mage check` to check the status of your local
development environment.

### Test

You can run the full test suite locally through Mage by running `mage test:ci`.
For these tests to pass you must have an instance of Postgres running that is
connectable and configured.

We're using [GitHub Actions](https://github.com/features/actions) to automate
our continus integration and delivery suite, which invokes `mage test:ci` to
determine the state of the build.

If your PostgreSQL database has a different password than `dev` you can pass a custom
database connection string like:

```shell
CAPE_DB_URL=postgres://postgres:<YOUR PASSWORD>@localhost:5432/cape mage test:ci
```

### Build

Once your development environment is bootstrapped you can build the Cape binary
by running the following command:

```
$ mage build:binary
```

This outputs the cape binary to `bin/cape`. To be able to run `cape` from anywhere
you can add it to your path like `export PATH=$PATH:<REPO ROOT>/cape/bin`. Add that
to your `.bashrc` to make it permenant.

### Deploy locally to Kubernetes

We're using [kind](https://kind.sigs.k8s.io/) for deploying Cape into Kubernetes locally.

Run the following `mage` command:

```
$ mage local:setup
```

This starts kind, builds the Docker images and deploys the Helm charts to get the coordinator and the database running.

If it completes successfully, you will see something like:

```
NAME: coordinator
LAST DEPLOYED: Fri May 15 16:23:30 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

to know it has completed successfully.

This creates an admin account for you with the email *cape_user@mycape.com*
and password *capecape*. If required, you can log in again with:

```bash
$ CAPE_PASSWORD=capecape cape login --email cape_user@mycape.com
```

### Troubleshooting

Here are some steps to try if you've installed Cape locally and you're having one of these problems:

#### Fedora: Docker does not support cgroups v2

Fedora 31 migrated from cgroups v1 to v2, but Docker and Kubernetes
don't support cgroups v2. Update the kernel to use cgroups v1.

```sh
sudo dnf install -y grubby && \
  sudo grubby \
  --update-kernel=ALL \
  --args="systemd.unified_cgroup_hierarchy=0"
```

#### Fedora: Unable to resolve host from inside the container

This is an
[issue](https://github.com/kubernetes-sigs/kind/issues/1547) that
surfaced in Fedora 32 when using kind.

Run `nmcli connection show` in your terminal to get the ethernet
interface device ID.

```sh
firewall-cmd --permanent --zone=trusted --add-interface=Docker0
firewall-cmd --get-zone-of-interface=<your eth interface>
firewall-cmd --zone=<zone from above> --add-masquerade --permanent
firewall-cmd --reload
```

Destroy your local environment with `mage local:destroy` and run
through the setup process once more.

## Submitting a pull request

To contribute, [fork](https://help.github.com/articles/fork-a-repo/) Cape, commit your changes, and [open a pull request](https://help.github.com/articles/using-pull-requests/).

While you may be asked to make changes to your submission during the review process, we will work with you on this and suggest changes. Consider giving us [push rights to your branch](https://help.github.com/articles/allowing-changes-to-a-pull-request-branch-created-from-a-fork/) so we can potentially also help via commits.

### Commit history and merging

For the sake of transparency our key rule is to keep a logical and intelligible commit history, meaning anyone stepping through the commits on either the `master` branch or as part of a review should be able to easily follow the changes made and their potential implications.

To this end we ask all contributors to sanitize pull requests before submitting them. All pull requests will either be [squashed or rebased](https://help.github.com/en/articles/about-pull-request-merges).

Some guidelines:

- Even simple code changes such as moving code around can obscure semantic changes, and in those case there should be two commits: for example, one that only moves code (with a note of this in the commit description) and one that performs the semantic change.

- Progressions that have no logical justification for being split into several commits should be squeezed.

- Code does not have to compile or pass all tests at each commit, but leave a remark and a plan in the commit description so reviewers are aware and can plan accordingly.

See below for some [useful tricks](#git-and-github) for working with Git and GitHub.

### Before submitting for review

Make sure to give some context and overview in the body of your pull request to make it easier for reviewers to understand your changes. Ideally explain why your particular changes were made the way they are.

Importantly, use [keywords](https://help.github.com/en/articles/closing-issues-using-keywords) such as `Closes #<issue-number>` to indicate any issues or other pull requests related to your work.

Furthermore:

- Run tests (`mage test:ci`) before submitting as our [CI](#continuous-integration) will block pull requests failing these checks.
- Test your change thoroughly with unit tests where appropriate.
- Update any affected comments in the code base.
- Add a line in [CHANGELOG.md](CHANGELOG.md) for any major change.

## Continuous integration

All pull requests are run against our [continuous integration suite](https://github.com/capeprivacy/cape/actions). The entire suite must pass before a pull request is accepted.

## Writing documentation

Ensure you add comments where necessary.

The documentation site is managed in the [documentation repository](https://github.com/capeprivacy/documentation).

## Useful tricks

### git and GitHub

- [GitHub Desktop](https://desktop.github.com/) provides a useful interface for inspecting and committing code changes
- `git add -p`
  - lets you leave out some changes in a file (GitHub Desktop can be used for this as well)
- `git commit --amend`
  - allows you to add to the previous commit instead of creating a new one
- `git rebase -i <commit>`
  - allows you to [squeeze and reorder commits](https://git-scm.com/book/en/v2/Git-Tools-Rewriting-History)
  - use `HEAD~5` to consider 5 most recent commits
  - use `<hash>~1` to start from commit identified by `<hash>`
- `git rebase master`
  - [pull in latest updates](https://git-scm.com/book/en/v2/Git-Branching-Rebasing) on `master`
- `git fetch --no-tags <repo> <remote branch>:<local branch>`
  - pulls down a remote branch from e.g. a fork and makes it available to check out as a local branch
  - `<repo>` is e.g. `git@github.com:<user>/tf-encrypted.git`
- `git push <repo> <local branch>:<remote branch>`
  - pushes the local branch to a remote branch on e.g. a fork
  - `<repo>` is e.g. `git@github.com:<user>/tf-encrypted.git`
- `git tag -d <tag> && git push origin :refs/tags/<tag>`
  - can be used to delete a tag remotely

## Reporting a bug

Please file [bug reports](https://github.com/capeprivacy/cape/issues/new?template=bug_report.md) as GitHub issues.

### Security disclosures

If you encounter a security issue then please responsibly disclose it by reaching out to us at [privacy@capeprivacy.com](privacy@capeprivacy.com). We will work with you to mitigate the issue and responsibly disclose it to anyone using the project in a timely manner.

## Asking for help

If you have any questions you are more than welcome to reach out through GitHub issues or [our Slack channel](https://join.slack.com/t/capecommunity/shared_invite/zt-f8jeskkm-r9_FD0o4LkuQqhJSa~~IQA).
