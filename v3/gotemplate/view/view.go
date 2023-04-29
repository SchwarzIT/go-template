package view

import "github.com/schwarzit/go-template/v3/gotemplate/module"

type EventType int

const (
	GenerateEvent EventType = iota
)

type Event struct {
	Type  EventType
	Value interface{}
}

type View interface {
	// Initialize the view with the available modules and a channel for communication with the engine.
	Start(modules []module.Module, events chan<- Event) error

	// Show a message to the user.
	ShowMessage(message string) error

	// Show an error message to the user.
	ShowError(err error) error
}
