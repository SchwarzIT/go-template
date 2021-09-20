package cmds

import (
	"fmt"
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/repos"
	"github.com/spf13/cobra"
)

const goTemplateGithubRepo = "https://github.com/schwarzit/go-template"

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

func checkVersion(_ *cobra.Command, _ []string) {
	tag, err := repos.LatestReleaseTag(goTemplateGithubRepo)
	if err != nil {
		printWarning("unable to fetch version information. There could be newer release for go/template.")
		return
	}

	if tag != config.Version {
		printWarning(fmt.Sprintf("newer version available: %s. Pls make sure to stay up to date to enjoy the latest features.", tag))
	}
}
