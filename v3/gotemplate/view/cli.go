package view

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/schwarzit/go-template/v3/gotemplate/module"
	"github.com/schwarzit/go-template/v3/gotemplate/option"
)

type CLI struct {
	State map[option.ModuleName]option.State
}

func NewCLI() *CLI {
	return &CLI{}
}

func (c *CLI) Start(modules []module.Module, events chan<- Event) error {
	// Print the welcome message.
	fmt.Println("Welcome to the Go Template Engine!")

	// Iterate over the modules and prompt the user for input.
	for _, m := range modules {
		name, err := m.GetName()
		if err != nil {
			return err
		}
		fmt.Printf("\n=== %s ===\n", name)

		options, err := m.GetOptions()
		if err != nil {
			return err
		}

		for _, o := range options {
			shouldDisplay, err := o.ShouldDisplay(nil)
			if err != nil {
				return err
			}
			if !shouldDisplay {
				continue
			}

			// Print the option title and description.
			title, err := o.GetTitle()
			if err != nil {
				return err
			}
			fmt.Println(title)

			desc, err := o.GetDescription()
			if err != nil {
				return err
			}
			fmt.Println(desc)

			defaultValue, err := o.GetDefaultValue(c.State)
			if err != nil {
				return err
			}
			if len(defaultValue) > 0 {
				fmt.Printf("Default value: %s\n", strings.Join(defaultValue, ", "))
			}

			// Print the available answers (if any).
			answers, err := o.GetAvailableAnswers()
			if err != nil {
				return err
			}
			if len(answers) > 0 {
				fmt.Printf("Available answers: %s\n", strings.Join(answers, ", "))
			}

			// Prompt the user for input.
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("> ")
			scanner.Scan()
			input := scanner.Text()

			// Validate the input and set the option value.
			if err := o.Validate([]string{input}); err != nil {
				fmt.Printf("Invalid input: %v\n", err)
				return err
			}
			if err := o.SetCurrentValue([]string{input}); err != nil {
				fmt.Printf("Failed to set option value: %v\n", err)
				return err
			}

			// Merge the state.
			moduleName, err := m.GetName()
			if err != nil {
				return err
			}
			if c.State == nil {
				c.State = make(map[option.ModuleName]option.State)
			}
			if c.State[moduleName] == nil {
				c.State[moduleName] = make(option.State)
			}
			c.State[moduleName][title], err = o.GetCurrentValue()
			if err != nil {
				return err
			}
		}
	}

	// Notify the engine that the user has confirmed the project generation.
	events <- Event{
		Type: GenerateEvent,
	}

	return nil
}

func (c *CLI) ShowMessage(message string) error {
	fmt.Println(message)
	return nil
}

func (c *CLI) ShowError(err error) error {
	fmt.Printf("Error: %v\n", err)
	return err
}
