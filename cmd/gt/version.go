package main

import (
	"fmt"

	"github.com/muesli/termenv"
	"github.com/schwarzit/go-template/pkg/colors"
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/spf13/cobra"
)

func buildVersionCommand(output *termenv.Output, gt *gotemplate.GT) *cobra.Command {
	goTemplateHighlighted := output.String(goTemplate).Foreground(output.Color(colors.Cyan))
	cmd := &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the version number of %s", goTemplateHighlighted),
		Long:  fmt.Sprintf("All software has versions. This is %s's.", goTemplateHighlighted),
		Run: func(cmd *cobra.Command, args []string) {
			gt.PrintVersion()
		},
	}

	return cmd
}
