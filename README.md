# `go/template`

[![SIT](https://img.shields.io/badge/SIT-awesome-blueviolet.svg)](https://jobs.schwarz)
[![CI](https://github.com/SchwarzIT/go-template/actions/workflows/main.yml/badge.svg)](https://github.com/SchwarzIT/go-template/actions/workflows/main.yml)
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

To get an overview of all options that can be set for the template you can take a look at the [options definition file](pkg/gotemplate/options.go).

## Usage

### Install

```bash
go install github.com/schwarzit/go-template/cmd/gt@latest
```

### Initialize your repo from the template

Use the template the generate your repo:

```bash
gt new
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
