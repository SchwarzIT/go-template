package option

// State represents the state of a module, as a map of string keys to string values.
type State map[string][]string

type ModuleName string
type Option interface {
	// GetTitle returns the title of the option.
	GetTitle() (string, error)

	// GetTemplateKey returns the template key of the option.
	GetTemplateKey() (string, error)

	// GetDescription returns the description of the option.
	GetDescription() (string, error)

	// GetDefaultValue returns the default value of the option, taking into account the
	// current state of the module.
	GetDefaultValue(state map[ModuleName]State) ([]string, error)

	// GetAvailableAnswers returns the available answers for the option. If the
	// option does not have a limited set of answers, it returns nil.
	GetAvailableAnswers() ([]string, error)

	// GetCurrentValue returns the current value of the option.
	GetCurrentValue() ([]string, error)

	// SetCurrentValue sets the value of the option.
	SetCurrentValue(value []string) error

	// Validate validates the given value for the option.
	Validate(value []string) error

	// ShouldDisplay determines if the option should be displayed given the
	// current state of the module.
	ShouldDisplay(state map[ModuleName]State) (bool, error)
}
