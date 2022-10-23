package gotemplate

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
	"github.com/schwarzit/go-template/pkg/colors"
)

func (gt *GT) colorStyler(color string) termenv.Style {
	return gt.styler().String().Foreground(gt.styler().Color(color))
}

func (gt *GT) cyanStyler() termenv.Style {
	return gt.colorStyler(colors.Cyan)
}

func (gt *GT) yellowStyler() termenv.Style {
	return gt.colorStyler(colors.Yellow)
}

func (gt *GT) printProgressf(format string, a ...interface{}) {
	s := gt.cyanStyler().Bold().Styled(fmt.Sprintf(format, a...))
	_, _ = fmt.Fprintln(gt.Out, s)
}

func (gt *GT) printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(gt.Out, format, a...)
}

func (gt *GT) printWarningf(format string, a ...interface{}) {
	warningBanner := gt.yellowStyler().Bold().Styled("WARNING")
	warningText := gt.yellowStyler().Styled(fmt.Sprintf(format, a...))

	_, _ = fmt.Fprintf(gt.Err, "%s: %s\n", warningBanner, warningText)
}

func (gt *GT) printOption(opts *Option, optionValues *OptionValues) {
	gt.printf("%s\n", gt.yellowStyler().Underline().Styled(opts.Description()))
	gt.printf("%s: (%v) ", gt.cyanStyler().Styled(opts.Name()), opts.Default(optionValues))
}

func (gt *GT) printBanner() {
	highlight := gt.cyanStyler().Styled
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
