package cmds

import (
	"fmt"
	"github.com/schwarzit/go-template/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of go/template.",
	Long:  "All software has versions. This is go/template's.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Version)
	},
}
