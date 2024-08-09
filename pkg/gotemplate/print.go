package gotemplate

import (
	"strings"

	"github.com/pterm/pterm"
)

func (gt *GT) printProgressf(format string, a ...interface{}) {
	pterm.NewStyle(pterm.FgCyan, pterm.Bold).Printfln(format, a...) // TODO: Use gt.Out
}

func (gt *GT) printf(format string, a ...interface{}) {
	pterm.Printf(format, a...) // TODO: Use gt.Out
}

func (gt *GT) printWarningf(format string, a ...interface{}) {
	pterm.Warning.Printf(format, a...) // TODO: Use gt.Out
}

func (gt *GT) printOption(opts *Option, optionValues *OptionValues) {
	gt.printf("%s\n", pterm.NewStyle(pterm.FgYellow, pterm.Underscore).Sprint(opts.Description()))
	gt.printf("%s: (%v) ", pterm.NewStyle(pterm.FgCyan).Sprint(opts.Name()), opts.Default(optionValues))
}

func (gt *GT) printBanner() {
	highlight := pterm.FgCyan.Sprint
	gt.printf("Hi! Welcome to the %s cli.\n", highlight("go/template"))
	gt.printf("This command will walk you through creating a new project.\n")
	gt.printf("You will first be asked to set values for the base paremeters that are needed for the minimal setup.\n")
	gt.printf("Afterwards you will get the opportunity to enable several extensions to extend the template's functionality.\n\n")
	gt.printf("Enter a value or leave blank to accept the (default), and press %s.\n", highlight("<ENTER>"))
	gt.printf("Press %s at any time to quit.\n\n", highlight("^C"))
}

func (gt *GT) printCategory(category string) {
	gt.printf(" --\n")
	gt.printf("| CATEGORY: %q\n", strings.ToUpper(category))
	gt.printf(" --\n")
}
