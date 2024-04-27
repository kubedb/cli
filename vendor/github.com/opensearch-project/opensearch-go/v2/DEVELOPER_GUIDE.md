- [Developer Guide](#developer-guide)
  - [Getting Started](#getting-started)
    - [Git Clone OpenSearch Go Client Repository](#git-clone-opensearch-go-client-repository)
    - [Install Prerequisites](#install-prerequisites)
      - [Go 1.11](#go-111)
      - [Docker](#docker)
      - [Windows](#windows)
    - [Unit Testing](#unit-testing)
    - [Integration Testing](#integration-testing)
      - [Execute integration tests from your terminal](#execute-integration-tests-from-your-terminal)
    - [Lint](#lint)
      - [Markdown lint](#markdown-lint)
  - [Use an Editor](#use-an-editor)
    - [GoLand](#goland)

# Developer Guide

So you want to contribute code to the OpenSearch Go Client? Excellent! We're glad you're here. Here's what you need to do:

## Getting Started

### Git Clone OpenSearch Go Client Repository

Fork [opensearch-project/opensearch-go](https://github.com/opensearch-project/opensearch-go) and clone locally, e.g. `git clone https://github.com/[your username]/opensearch-go.git`.

### Install Prerequisites

#### Go 1.11

OpenSearch Go Client builds using [Go](https://golang.org/doc/install) 1.11 at a minimum.

#### Docker

[Docker](https://docs.docker.com/install/) is required for building some OpenSearch artifacts and executing integration tests.

#### Windows

To build the project on Windows, use [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install), the compatibility layer for running Linux applications.

Install ```make```
```
sudo apt install make
```

### Unit Testing

Go has a simple tool for running tests, and we simplified it further by creating this make command:

```
make test-unit
```

Individual unit tests can be run with the following command:

```
cd folder-path/to/test;
go test -v -run TestName;
```

### Integration Testing

In order to test opensearch-go client, you need a running OpenSearch cluster. You can use Docker to accomplish this. The [Docker Compose file](.ci/opensearch/docker-compose.yml) supports the ability to run integration tests for the project in local environments. If you have not installed docker-compose, you can install it from this [link](https://docs.docker.com/compose/install/).

In order to differentiate unit tests from integration tests, Go has a built-in mechanism for allowing you to logically separate your tests with [build tags](https://pkg.go.dev/cmd/go#hdr-Build_constraints). The build tag needs to be placed as close to the top of the file as possible, and must have a blank line beneath it. Hence, create all integration tests with build tag 'integration'.

#### Execute integration tests from your terminal

1. Run below command to start containers. By default, it will launch latest OpenSearch cluster.
   ```
   make cluster.build cluster.start
   ```
2. Run all integration tests.
   ```
   make test-integ race=true
   ```
3. Stop and clean containers.
   ```
   make cluster.stop cluster.clean
   ```

## Lint

To keep all the code in a certain uniform format, it was decided to use some writing rules. If you wrote something wrong, it's okay, you can simply run the script to check the necessary files, and optionally format the content. But keep in mind that all these checks are repeated on the pipeline, so it's better to check locally.

### Markdown lint

To check the markdown files, run the following command:

```
make lint.markdown
```

## Use an Editor

### GoLand

You can import the OpenSearch project into GoLand as follows:

1. Select **File | Open**
2. In the subsequent dialog navigate to the ~/go/src/opensearch-go and click **Open**

After you have opened your project, you need to specify the location of the Go SDK. You can either specify a local path to the SDK or download it. To set the Go SDK, navigate to **Go | GOROOT** and set accordingly.
