package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/spf13/cobra"
)

func buildNewCommand(gt *gotemplate.GT) *cobra.Command {
	var (
		configFile string
		opts       gotemplate.NewRepositoryOptions
	)

	underline := color.New(color.Underline).SprintFunc()

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project repository",
		Long: fmt.Sprintf(`Create a new Golang project folder using the "_template" folder in github.com/schwarzit/go-template as base.

Since this is only a template some parameters are needed to render the final project folder.
This saves you time since you don't need to find+replace anymore.

There are two available modes to set those parameters:

%s
By default the CLI will run in Interactive Mode.
This means all parameters values will be gathered through stdin user input.
To use that just type plain "gt new" and follow the further instructions.

%s
Since interactive user input is not a feasible solution in all cases there's also the option
to a pass config file through the "--config" flag.
This defines the parameters as key value pairs.
To get further information look at the flag's documentation.
`, underline("Interactive Mode"), underline("File Mode")),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			configValues, err := getValues(gt, configFile)
			if err != nil {
				return err
			}

			opts.OptionValues = configValues
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return gt.InitNewProject(&opts)
		},
	}

	cmd.Flags().StringVarP(
		&configFile,
		"config", "c", "",
		`Config file that defines all parameters.
This is helpful if you don't want to run the CLI interactively.
It should either be a json or a yaml file.
An example file could look like:

// values.yaml
parameters:
    projectName: Some Project
    projectSlug: some-project
    projectDescription: Some random project
    appName: somecli
    moduleName: github.com/some-user/some-project
    golangciVersion: 1.42.1
integrations:
    grpcEnabled: true
    grpcGatewayEnabled: false`,
	)

	cmd.Flags().StringVar(
		&opts.CWD, "cwd", "./",
		`Current working directory.
Can be set to decide where to create the new project folder,
`)

	return cmd
}

func getValues(gt *gotemplate.GT, configFile string) (*gotemplate.OptionValues, error) {
	if configFile != "" {
		return gt.LoadConfigValuesFromFile(configFile)
	}
	return gt.LoadConfigValuesInteractively()
}
