package gotemplate

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func (gt *GT) printProgressf(format string, a ...interface{}) {
	_, _ = color.New(color.FgCyan, color.Bold).Fprintf(gt.Out, format, a...)
	_, _ = fmt.Fprintln(gt.Out)
}

func (gt *GT) printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(gt.Out, format, a...)
}

func (gt *GT) printWarningf(format string, a ...interface{}) {
	headerHighlight := color.New(color.FgYellow, color.Bold).SprintFunc()
	highlight := color.New(color.FgYellow)

	_, _ = fmt.Fprintf(gt.Err, "%s: ", headerHighlight("WARNING"))
	_, _ = highlight.Fprintf(gt.Err, format, a...)
	_, _ = fmt.Fprintln(gt.Err)
}

func (gt *GT) printOption(opts *Option, optionValues *OptionValues) {
	highlight := color.New(color.FgCyan).SprintFunc()
	underline := color.New(color.FgHiYellow, color.Underline).SprintFunc()
	gt.printf("%s\n", underline(opts.Description()))
	gt.printf("%s: (%v) ", highlight(opts.Name()), opts.Default(optionValues))
}

func (gt *GT) printBanner() {
	highlight := color.New(color.FgCyan).SprintFunc()
	gt.printf("Hi! Welcome to the %s cli.\n", highlight("go/template"))
	gt.printf("This command will walk you through creating a new project.\n")
	gt.printf("You will first be asked to set values for the base paremeters that are needed for the minimal setup.\n")
	gt.printf("Afterwards you will get the opportunity to enable several extensions to extend the template's functionality.\n\n")
	gt.printf("Enter a value or leave blank to accept the (default), and press %s.\n", highlight("<ENTER>"))
	gt.printf("Press %s at any time to quit.\n\n", highlight("^C"))
}

func (gt *GT) printCategory(category string) {
	gt.printf(" --\n")
	gt.printf("| CATEGORY: %q\n", strings.ToUpper(category))
	gt.printf(" --\n")
}
