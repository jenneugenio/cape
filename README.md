# Cape

![](https://github.com/capeprivacy/cape/workflows/Main/badge.svg)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Gitub release]("https://img.shields.io/github/v/release/capeprivacy/cape.svg?logo=github")](https://github.com/capeprivacy/cape/releases/latest)
[![Chat on Slack](https://img.shields.io/badge/chat-on%20slack-7A5979.svg)](https://join.slack.com/t/capecommunity/shared_invite/zt-f8jeskkm-r9_FD0o4LkuQqhJSa~~IQA)

Cape Core contains the key functionality of Cape Privacy, including:

* A CLI (command line interface)
* Cape Coordinator, which provides policy management workflows, and controllers to work with the Cape Privacy libraries.

View the [documentation](https://docs.capeprivacy.com/cape-core/).

## Getting started

Cape Core comprises multiple elements, some of which interact with other parts of the Cape ecosystem. The following links point to the correct getting started resources for each element:

* [Cape Coordinator installation guides](https://docs.capeprivacy.com/cape-core/coordinator)
* [Cape CLI installation guide](https://docs.capeprivacy.com/cape-core/cli/installation/)
* [Cape CLI usage guide](https://docs.capeprivacy.com/cape-core/cli/usage/)
* [Tutorial: using Coordinator and CLI with Cape Python](https://docs.capeprivacy.com/libraries/cape-python/coordinator-quickstart/)

For information on how to work on Cape development locally, refer to [Contributing](./CONTRIBUTING.md).


## About Cape Privacy and Cape

[Cape Privacy](https://capeprivacy.com) helps teams share data and make decisions for safer and more powerful data science. Learn more at [capeprivacy.com](https://capeprivacy.com).

Cape contains the core functionality of Cape Privacy, including a CLI (command line interface), policy management workflow, and controllers to work with the [Cape Privacy libraries](https://docs.capeprivacy.com/libraries/).

### Cape architecture

Cape is comprised of multiples services and libraries. The Coordinator provides policy and user management through a CLI. This in turn can interact with Cape's libraries, applying your data policy directly to the data transformation scripts.

### Project status and roadmap

Cape Core was released 30th July 2020. It is actively maintained and developed, alongside other elements of the Cape ecosystem.

**Upcoming features:**

* Audit logging configuration: set up configuration for how and where you log actions in Cape Coordinator, such as project and policy creation, user changes and user actions in Cape.
* Governance tooling: integrate basic data governance information to be used within Cape Coordinator for writing better policy, with a possible integration with Apache Atlas or other open-source governance tools.
* Pipeline orchestrator integration: ability to connect with Spark orchestration tools (such as YARN, Mesos, and Airflow) and pull information on jobs that are running for easier management of running Spark installations.

The goal is a complete data management ecosystem. Cape Privacy provides [Cape Coordinator](https://docs.capeprivacy.com/cape-core/coordinator/), to manage policy and users. This will interact with the Cape Privacy libraries (such as [Cape Python](https://docs.capeprivacy.com/libraries/cape-python/)) through a workers interface, and with your own data services through an API.

## Help and resources

If you need help using Cape, you can:

* View the [documentation](https://docs.capeprivacy.com/).
* Submit an issue.
* Talk to us on our [community Slack](https://join.slack.com/t/capecommunity/shared_invite/zt-f8jeskkm-r9_FD0o4LkuQqhJSa~~IQA).

Please file [feature requests](https://github.com/capeprivacy/cape/issues/new?template=feature_request.md) and
[bug reports](https://github.com/capeprivacy/cape/issues/new?template=bug_report.md) as GitHub issues.

## Community

[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/0)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/0)[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/1)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/1)[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/2)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/2)[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/3)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/3)[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/4)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/4)[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/5)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/5)[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/6)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/6)[![](https://sourcerer.io/fame/justin1121/capeprivacy/cape/images/7)](https://sourcerer.io/fame/justin1121/capeprivacy/cape/links/7)

### Contributing

View our [contributing](CONTRIBUTING.md) guide for more information.

### Code of conduct

Our [code of conduct](https://capeprivacy.com/conduct/) is included on the Cape Privacy website. All community members are expected to follow it. Please refer to that page for information on how to report problems.


## License

Licensed under Apache License, Version 2.0 (see [LICENSE](https://github.com/capeprivacy/cape/blob/master/LICENSE) or http://www.apache.org/licenses/LICENSE-2.0). Copyright as specified in [NOTICE](https://github.com/capeprivacy/cape/blob/master/NOTICE).