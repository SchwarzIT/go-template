package gotemplate

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/schwarzit/go-template/pkg/option"
)

func (gt *GT) printProgress(str string) {
	_, _ = color.New(color.FgCyan, color.Bold).Fprintln(gt.Out, str)
}

func (gt *GT) printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(gt.Out, format, a...)
}

func (gt *GT) printWarning(str string) {
	headerHighlight := color.New(color.FgYellow, color.Bold).SprintFunc()
	highlight := color.New(color.FgYellow).SprintFunc()

	_, _ = fmt.Fprintf(gt.Err, "%s: %s\n", headerHighlight("WARNING"), highlight(str))
}

func (gt *GT) printOption(opts *option.Option) {
	highlight := color.New(color.FgCyan).SprintFunc()
	underline := color.New(color.FgHiYellow, color.Underline).SprintFunc()
	gt.printf("%s\n", underline(opts.Description))
	gt.printf("%s: (%v) ", highlight(opts.Name), opts.Default)
}

func (gt *GT) printBanner() {
	highlight := color.New(color.FgCyan).SprintFunc()
	gt.printf("Hi! Welcome to the %s cli.\n", highlight("go/template"))
	gt.printf("This command will walk you through creating a new project.\n\n")
	gt.printf("Enter a value or leave blank to accept the (default), and press %s.\n", highlight("<ENTER>"))
	gt.printf("Press %s at any time to quit.\n\n", highlight("^C"))
}
