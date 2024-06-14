# `saga`

 
[![Go Report Card](https://goreportcard.com/badge/github.com/legendaryum-metaverse/saga)](https://goreportcard.com/report/github.com/legendaryum-metaverse/saga)
[![Go Reference](https://pkg.go.dev/badge/github.com/legendaryum-metaverse/saga?status.svg)](https://pkg.go.dev/github.com/legendaryum-metaverse/saga?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/legendaryum-metaverse/saga/-/badge.svg)](https://sourcegraph.com/github.com/legendaryum-metaverse/saga?badge)
[![Release](https://img.shields.io/github/release/legendaryum-metaverse/saga.svg?style=flat-square)](https://github.com/legendaryum-metaverse/saga/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/legendaryum-metaverse/saga/release.yml?branch=main)](https://github.com/legendaryum-metaverse/saga/tree/main)
[![GitHub](https://img.shields.io/github/license/legendaryum-metaverse/saga)](https://github.com/legendaryum-metaverse/saga/blob/main/LICENSE)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/legendaryum-metaverse/saga)](https://github.com/legendaryum-metaverse/saga/pulse)
[![GitHub last commit](https://img.shields.io/github/last-commit/legendaryum-metaverse/saga)](https://github.com/legendaryum-metaverse/saga/commits/main)

[**saga**](https://pkg.go.dev/github.com/legendaryum-metaverse/saga) is a Go library designed to streamline communication
between microservices using RabbitMQ. It enables easy implementation of event-driven architectures and saga patterns, while ensuring reliable message delivery.

## Features

**Core Communication:**

-   **Publish/Subscribe Messaging:** Exchange messages between microservices using a
    publish-subscribe pattern.
-   **Headers-Based Routing:** Leverage the power of RabbitMQ's headers exchange for flexible and dynamic routing of messages based on custom headers.
-   **Durable Exchanges and Queues:** Ensure message persistence and reliability with durable RabbitMQ components.

**Saga Management:**

<div style="text-align: center;">
<img src="https://raw.githubusercontent.com/legendaryum-metaverse/legend-transactional/main/.github/assets/saga.png" alt="legendaryum" style="width: 90%;"/>
</div>

-   **Saga Orchestration:** Coordinate complex, multi-step transactions across multiple microservices with saga orchestration.
-   **Saga Step Handlers:** Implement step-by-step saga logic in your microservices using callbacks.
-   **Compensation Logic:** Define compensating actions for saga steps to handle failures
    gracefully and maintain data consistency.

## Contributors

Thanks to [all contributors](https://github.com/legendaryum-metaverse/saga/graphs/contributors)!

## Author

Jorge Clavijo <https://github.com/jym272>

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.
