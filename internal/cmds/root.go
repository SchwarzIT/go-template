package cmds

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:              "go-template",
	Short:            "go/template is a tool for jumpstarting production-ready Golang projects quickly",
	Long:             "A repo template generator build by schwarzit. Full documentation at github.com/SchwarzIT/go-template.",
	PersistentPreRun: checkVersion,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printError(err)
		os.Exit(1)
	}
}
