# 2. Build own CLI

Date: 2021-10-11

## Status

Accepted

## Context

At the beginning of the project [cookiecutter](https://github.com/cookiecutter/cookiecutter) was used to generate and parameterize the template folder.
The following things have been considered in the Discussion:

`cookiecutter` is a ready to use tool which made it perfect to get up and running quickly with the template.
When the decision was made that the project should be open-sourced to a possibly larger audience this was rethought.
The first problem that occured with `cookiecutter` is the dependency on python and the usage of python to write custom logic.
One major advantage of Go are the statically compiled binaries that run everywhere. That's why it felt wrong to include a dependency in a Go related tool.
Also for maintaining and contributing to the template using another language is questionable.

The major advantage of a self developed CLI on the other hand is maximum freedom when choosing features.
With that Go's `text/template` can be used over the more Python specific [Jinja templates](https://jinja.palletsprojects.com/en/3.0.x/) used in `cookiecutter`.
Also the structure of the options and way of collecting user input can be refined.

## Decision

Cookiecutter should be replaced with a self developed Go CLI.

## Consequences

- Project structure is changed to use `_template` folder instead of `{{cookiecutter.project_slug}}` to keep the folder at the top and visible
- CLI is developed using Go standard library's `text/template`
- CLI needs to distributed in all common package managers
