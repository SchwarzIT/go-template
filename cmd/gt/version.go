package main

import (
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/spf13/cobra"
)

func buildVersionCommand(gt *gotemplate.GT) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of go/template.",
		Long:  "All software has versions. This is go/template's.",
		Run: func(cmd *cobra.Command, args []string) {
			gt.PrintVersion()
		},
	}

	return cmd
}
