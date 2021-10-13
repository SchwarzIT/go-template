# 3. Use Go structs to define the available options

Date: 2021-10-11

## Status

Accepted

## Context

When creating the Go CLI for the template parameterization a YAML file was used to define all available options/paramaters to be used in the template.
This was inspired by `cookiecutter`, the previously used tool.

While extending the options further custumization options were added to this YAML file like for example regex validators for each option or a field to define files that should be removed/ added whenever an option is true or false. The complexity in the YAML file/ the custom syntax used quickly evolved into its own pseudo language with conditions and a growing set of options to customize the parameters with.

At this point the benefits of the YAML file was questioned since the only reason to have it was basically the `cookiecutter` origins of the repo.

Go structs on the other hand have the advantage that any custom logic and conditions can be included much easier since it's a fully featured programming language as opposed to YAML.
Also contributing should be easier with Go structs to define options since there's no need to learn all the options available in the YAML DSL.

## Decision

Use Go structs instead of YAML to define the available options/parameters and their behaviour (default values, dependencies etc.)

## Consequences

- YAML file is replaced with Go structs
- Generalized Go struct properties to maximize extensibility
