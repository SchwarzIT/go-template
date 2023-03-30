package option

// Option represents a configuration option that can be set by the user.
type Option interface {
	// GetTitle returns the title of the option.
	GetTitle() string

	// GetDescription returns the description of the option.
	GetDescription() string

	// GetDefaultValue returns the default value of the option.
	GetDefaultValue() string

	// GetAvailableAnswers returns a list of available answers for the option.
	GetAvailableAnswers() []string

	// GetValue returns the current value of the option.
	GetValue() string

	// SetValue sets the value of the option to the given value.
	// Returns an error if the value is invalid.
	SetValue(value string) error

	// Validate validates the given value for the option.
	// Returns an error if the value is invalid.
	Validate(value string) error
}
