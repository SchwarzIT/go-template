package main

import (
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/spf13/cobra"
)

func buildNewCommand(gt *gotemplate.GT) *cobra.Command {
	var (
		configFile string
		opts       gotemplate.NewRepositoryOptions
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project repository.",
		Long:  "Fill out all given parameters to configure and jump start your next project repository.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			configValues, err := getValues(gt, configFile)
			if err != nil {
				return err
			}

			opts.ConfigValues = configValues
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
It should either be a json or a yaml file.`,
	)

	cmd.Flags().StringVar(
		&opts.CWD, "cwd", "./",
		`Current working directory.
Can be set to decide where to create the new project folder.,
`)

	return cmd
}

func getValues(gt *gotemplate.GT, configFile string) (map[string]interface{}, error) {
	if configFile != "" {
		return gt.LoadConfigValuesFromFile(configFile)
	}
	return gt.LoadConfigValuesInteractively()
}
