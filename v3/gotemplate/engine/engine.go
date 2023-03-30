package engine

import (
	"github.com/schwarzit/go-template/v3/gotemplate/module"
	"github.com/schwarzit/go-template/v3/gotemplate/view"
)

type Engine struct {
	modules []module.Module
}

func NewEngine(modules []module.Module) *Engine {
	return &Engine{
		modules: modules,
	}
}

func (e *Engine) Start(v view.View) error {
	// Create a channel for communication with the view.
	events := make(chan view.Event)

	// Listen for events from the view and generate the project when necessary.
	go func() {
		for event := range events {
			if event.Type == view.GenerateEvent {
				if err := e.generateProject(); err != nil {
					v.ShowError(err)
				} else {
					v.ShowMessage("Project generated successfully.")
				}
			}
		}
	}()

	// Initialize the view.
	if err := v.Start(e.modules, events); err != nil {
		return err
	}

	return nil
}

func (e *Engine) generateProject() error {
	// Generate the output for each module using the current option values.
	for _, m := range e.modules {
		if err := m.Generate(); err != nil {
			return err
		}
	}

	return nil
}
