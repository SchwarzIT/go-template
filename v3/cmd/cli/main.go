package main

import (
	"fmt"
	"os"

	"github.com/schwarzit/go-template/v3/gotemplate/engine"
	"github.com/schwarzit/go-template/v3/gotemplate/module"
	"github.com/schwarzit/go-template/v3/gotemplate/view"
	"github.com/schwarzit/go-template/v3/modules/base"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Create a new module with two options: project name and description.
	m, err := base.NewReadmeModule()
	if err != nil {
		return err
	}

	// Create a new CLI view.
	cli := view.NewCLI()

	// Create a new engine with the module and the CLI view.
	engine := engine.NewEngine([]module.Module{m})

	// Start the engine.
	err = engine.Start(cli)
	if err != nil {
		return err
	}

	return nil
}
