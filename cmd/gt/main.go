package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func printError(err error) {
	headerHighlight := color.New(color.FgRed, color.Bold).SprintFunc()
	highlight := color.New(color.FgRed).SprintFunc()

	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", headerHighlight("ERROR"), highlight(err.Error()))
}

func buildRootCommand() *cobra.Command {
	gt := gotemplate.New()

	cmd := &cobra.Command{
		Use:   "gt",
		Short: "gt is go/template's cli for jumpstarting production-ready Golang projects quickly",
		Long:  "A repo template generator build by schwarzit. Full documentation at github.com/SchwarzIT/go-template.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Enable swapping out stdout/stderr for testing
			gt.Out = cmd.OutOrStdout()
			gt.Err = cmd.OutOrStderr()
			gt.In = cmd.InOrStdin()

			gt.CheckVersion()
		},
		// don't show errors and usage on errors in any RunE function.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(buildNewCommand(gt))
	cmd.AddCommand(buildVersionCommand(gt))

	return cmd
}
