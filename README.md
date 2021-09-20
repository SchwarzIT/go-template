# `go/template`

[![SIT](https://img.shields.io/badge/SIT-awesome-blueviolet.svg)](https://jobs.schwarz)
[![Semgrep](https://github.com/SchwarzIT/go-template/actions/workflows/semgrep.yml/badge.svg)](https://github.com/SchwarzIT/go-template/actions/workflows/semgrep.yml)

`go/template` is a tool for jumpstarting production-ready Golang projects quickly.

## Batteries included

- Makefile for most common tasks
- optimized Dockerfile
- golangci-lint default configuration
- pre-push git hook to ensure no linting issues
- gRPC support
- folder structure based on [github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- enforced default packages
  - `go.uber.org/zap` for logging
  - `go.uber.org/automaxprocs` to be safe in container environments (see [this article](https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs) for more information)
  
## Options

| Option                 | Description                                                                                                                                                                                                                                                                                                 |
| :--------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `project_name`         | The name of your project. This will also end up in the `README.md` and will be converted to be a valid default folder name (`project_slug`).                                                                                                                                                                |
| `project_slug`         | The name for the newly created project folder.                                                                                                                                                                                                                                                              |
| `project_description`  | The description of the project. This will be used in the `README.md`.                                                                                                                                                                                                                                       |
| `app_name`             | The name of the binary that you want to create. Could be the same your `project_slug` but since Go supports multiple apps in one repo it could also be sth. else. For example if your project is for some API there could be one app for the server and one CLI client.                                     |
| `module_name`          | The name of the Go module defined in the `go.mod` file. This is used if you want to `go get` the module. Please be aware that this depends on your version control system. The default points to `github.com` but for devops for example it would look sth. like this: `dev.azure.com/org/project/repo.git` |
| `golangci_version`     | The version of `golangci-lint` that you'd like to use for linting.                                                                                                                                                                                                                                          |
| `grpc_enabled`         | If enabled the created project will contain an example protobuf definition as well as several tools needed for gRPC development like `buf`. Also the Makefile will contain more targets to support the workflow.                                                                                            |
| `grpc_gateway_enabled` | If enabled the required dependencies for the [grpc-gateway project](https://github.com/grpc-ecosystem/grpc-gateway) will be added.                                                                                                                                                                          |

## Usage

### Install

```bash
go get -u github.com/SchwarzIT/go-template
```

### Initialize your repo from the template

Use the template the generate your repo:

```bash
cookiecutter https://bitbucket.schwarz/scm/goc/template.git
```

Initialize the project:

```bash
cd <your project>
make all
```

## Maintainers

| Name                                           | Email                        |
| :--------------------------------------------- | :--------------------------- |
| [@brumhard](https://github.com/brumhard)       | tobias.brumhard@mail.schwarz |
| [@linuxluigi](https://github.com/linuxluigi)   | steffen.exler@mail.schwarz   |
| [@danielzwink](https://github.com/danielzwink) | daniel.zwink@mail.schwarz    |

## Contribution

Contributions are very much appreciated.  
If you have anything to add to the template you are welcome to open a PR.
If your idea contains some major changes please open an issue to discuss first.
