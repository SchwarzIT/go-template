package view

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/schwarzit/go-template/v3/gotemplate/module"
)

type CLI struct{}

func NewCLI() *CLI {
	return &CLI{}
}

func (c *CLI) Start(modules []module.Module, events chan<- Event) error {
	// Print the welcome message.
	fmt.Println("Welcome to the Go Template Engine!")

	// Iterate over the modules and prompt the user for input.
	for _, m := range modules {
		fmt.Printf("\n=== %s ===\n", m.GetName())

		for _, o := range m.GetOptions() {
			// Print the option title and description.
			fmt.Println(o.GetTitle())
			fmt.Println(o.GetDescription())

			// Print the available answers (if any).
			answers := o.GetAvailableAnswers()
			if len(answers) > 0 {
				fmt.Printf("Available answers: %s\n", strings.Join(answers, ", "))
			}

			// Prompt the user for input.
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("> ")
			scanner.Scan()
			input := scanner.Text()

			// Validate the input and set the option value.
			if err := o.Validate(input); err != nil {
				fmt.Printf("Invalid input: %v\n", err)
				return err
			}
			if err := o.SetValue(input); err != nil {
				fmt.Printf("Failed to set option value: %v\n", err)
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

func (c *CLI) ShowMessage(message string) {
	fmt.Println(message)
}

func (c *CLI) ShowError(err error) {
	fmt.Printf("Error: %v\n", err)
}
