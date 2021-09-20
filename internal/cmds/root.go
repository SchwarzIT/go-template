package cmds

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "go-template",
	Short: "go/template is a tool for jumpstarting production-ready Golang projects quickly",
	Long:  "A repo template generator build by schwarzit. Full documentation at github.com/SchwarzIT/go-template.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func printError(err error) {
	headerHighlight := color.New(color.FgRed, color.Bold).SprintFunc()
	highlight := color.New(color.FgRed).SprintFunc()

	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", headerHighlight("Error during execution"), highlight(err.Error()))
}
