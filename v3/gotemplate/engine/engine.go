package engine

import (
	"sync"

	"github.com/schwarzit/go-template/v3/gotemplate/module"
	"github.com/schwarzit/go-template/v3/gotemplate/option"
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

	var wg sync.WaitGroup
	wg.Add(1)

	// Listen for events from the view and generate the project when necessary.
	go func() {
		defer wg.Done()
		for event := range events {
			if event.Type == view.GenerateEvent {
				if err := e.generateProject(); err != nil {
					v.ShowError(err)
				} else {
					v.ShowMessage("Project generated successfully.")
				}
				return
			}
		}
	}()

	// Initialize the view.
	if err := v.Start(e.modules, events); err != nil {
		return err
	}

	// Wait for the event listener goroutine to finish.
	wg.Wait()

	return nil
}

func (e *Engine) generateProject() error {
	// Collect all options from the modules.
	var allOptions []option.Option
	for _, m := range e.modules {
		options, err := m.GetOptions()
		if err != nil {
			return err
		}
		allOptions = append(allOptions, options...)
	}

	// Generate the output for each module using the current option values.
	for _, m := range e.modules {
		templateFiles, err := m.GetTemplateFiles()
		if err != nil {
			return err
		}

		if err := m.Generate(templateFiles, allOptions); err != nil {
			return err
		}
	}

	return nil
}
