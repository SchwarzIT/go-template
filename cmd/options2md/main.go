package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"text/template"

	"github.com/schwarzit/go-template/pkg/gotemplate"
)

// options2md tranlates all options available in go/template to a markdown file
// defined by the -o flag. This can be used for documentation.
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var outputFile string
	flag.StringVar(&outputFile, "o", "./options.md", "The file to write")
	flag.CommandLine.Parse(args)

	if outputFile == "" {
		return errors.New("`o` is a required parameter")
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, gotemplate.NewOptions(nil))
}

var (
	tmplString = `<!-- Code generated by options2md. DO NOT EDIT. -->
# Options

The following sections describe all options that are currently available for go/template for templating.
The options are divided into base options and extension options.
Base options are needed for the minimal base template and are mandatory in any case.
The extension options on the other hand enable optional features in the template such as gRPC support or open source lincenses.

## Base

| Name | Description |
| :--- | :---------- |
{{- range $index, $option := .Base}}
| {{ $option.Name | code }} | {{ $option.Description | replace "\n" "<br>" }} |
{{- end}}

## Extensions
{{- range $index, $category := .Extensions}}

### {{ $category.Name | code }}

| Name | Description |
| :--- | :---------- |
{{- range $index, $option := $category.Options}}
| {{ $option.Name | code }} | {{ $option.Description | replace "\n" "<br>" }} |
{{- end}}
{{- end}}
`
	funcMap = template.FuncMap{
		"replace": func(old, new, src string) string {
			return strings.ReplaceAll(src, old, new)
		},
		"code": func(s string) string {
			return fmt.Sprintf("`%s`", s)
		},
	}
	tmpl = template.Must(template.New("").Funcs(funcMap).Parse(tmplString))
)
