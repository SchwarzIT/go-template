package gotemplate

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/schwarzit/go-template/pkg/option"
)

func (gt *GT) printProgress(str string) {
	_, _ = color.New(color.FgCyan, color.Bold).Fprintln(gt.Out, str)
}

func (gt *GT) printWarning(str string) {
	headerHighlight := color.New(color.FgYellow, color.Bold).SprintFunc()
	highlight := color.New(color.FgYellow).SprintFunc()

	_, _ = fmt.Fprintf(gt.Err, "%s: %s\n", headerHighlight("WARNING"), highlight(str))
}

func (gt *GT) printOption(opts *option.Option) {
	highlight := color.New(color.FgCyan).SprintFunc()
	underline := color.New(color.FgHiYellow, color.Underline).SprintFunc()
	_, _ = fmt.Fprintf(gt.Out, "%s\n", underline(opts.Description))
	_, _ = fmt.Fprintf(gt.Out, "%s: (%v) ", highlight(opts.Name), opts.Default)
}

func (gt *GT) printBanner() {
	highlight := color.New(color.FgCyan).SprintFunc()
	_, _ = fmt.Fprintf(gt.Out, "Hi! Welcome to the %s cli.\n", highlight("go/template"))
	_, _ = fmt.Fprintf(gt.Out, "This command will walk you through creating a new project.\n\n")
	_, _ = fmt.Fprintf(gt.Out, "Enter a value or leave blank to accept the (default), and press %s.\n", highlight("<ENTER>"))
	_, _ = fmt.Fprintf(gt.Out, "Press %s at any time to quit.\n\n", highlight("^C"))
}
