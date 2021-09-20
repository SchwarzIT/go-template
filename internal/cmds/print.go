package cmds

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func printProgress(str string) {
	_, _ = color.New(color.FgCyan, color.Bold).Println(str)
}

func printError(err error) {
	headerHighlight := color.New(color.FgRed, color.Bold).SprintFunc()
	highlight := color.New(color.FgRed).SprintFunc()

	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", headerHighlight("ERROR"), highlight(err.Error()))
}

func printWarning(str string) {
	headerHighlight := color.New(color.FgYellow, color.Bold).SprintFunc()
	highlight := color.New(color.FgYellow).SprintFunc()

	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", headerHighlight("WARNING"), highlight(str))
}
