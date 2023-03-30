package main

import (
	"github.com/schwarzit/go-template/v3/gotemplate/engine"
	"github.com/schwarzit/go-template/v3/gotemplate/module"
	"github.com/schwarzit/go-template/v3/gotemplate/option"
	"github.com/schwarzit/go-template/v3/gotemplate/view"
)

func main() {
	// Create a new module with two options: project name and description.
	m := module.NewBaseModule(
		"Readme",
		[]option.Option{
			option.NewBaseOption(
				"Project Name",
				"The name of the new project.",
				"My Project",
				nil,
				nil,
			),
			option.NewBaseOption(
				"Project Description",
				"A brief description of the new project.",
				"This is a new project.",
				nil,
				nil,
			),
		},
	)

	// Create a new CLI view.
	cli := view.NewCLI()

	// Create a new engine with the module and the CLI view.
	engine := engine.NewEngine([]module.Module{m})

	// Start the engine.
	if err := engine.Start(cli); err != nil {
		panic(err)
	}
}
