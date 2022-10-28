package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/muesli/termenv"
	"github.com/schwarzit/go-template/pkg/colors"
	"github.com/schwarzit/go-template/pkg/gotemplate"
	"github.com/spf13/cobra"
)

const goTemplate = "go/template"

func main() {
	output := termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.EnvColorProfile()))

	cmd := buildRootCommand(output)
	if err := cmd.Execute(); err != nil {
		printError(output, err)
		os.Exit(1)
	}
}

func printError(output *termenv.Output, err error) {
	redStyler := output.String().Foreground(output.Color(colors.Red))

	_, _ = fmt.Fprintf(
		os.Stderr, "%s: %s\n",
		redStyler.Bold().Styled("ERROR"),
		redStyler.Styled(err.Error()),
	)
}

func buildRootCommand(output *termenv.Output) *cobra.Command {
	gt := gotemplate.New()

	cmd := &cobra.Command{
		Use:   "gt",
		Short: "gt is go/template's cli for jumpstarting production-ready Golang projects quickly",
		Long: fmt.Sprintf(`%[1]s is a tool for jumpstarting production-ready Golang projects quickly.

To begin working with %[1]s, run the 'gt new' command:

	$ gt new

This will prompt you to create a new Golang project using standard configs.

For more information, please visit the project's Github page: github.com/schwarzit/go-template.`,
			output.String(goTemplate).Foreground(output.Color(colors.Cyan)),
		),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Enable swapping out stdout/stderr for testing
			gt.Out = cmd.OutOrStdout()
			gt.Err = cmd.OutOrStderr()
			gt.InScanner = bufio.NewScanner(cmd.InOrStdin())

			gt.CheckVersion()
		},
		// don't show errors and usage on errors in any RunE function.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(buildNewCommand(output, gt))
	cmd.AddCommand(buildVersionCommand(output, gt))

	return cmd
}
