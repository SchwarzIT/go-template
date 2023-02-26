package gotemplate

// View is an interface that defines the methods needed for a user interface.
type View interface {
	// PresentQuestion presents a question to the user and returns the user's response.
	PresentQuestion(question TemplateQuestion) (*TemplateQuestion, error)

	// // ShowMessage displays a message to the user.
	// ShowMessage(message string) error

	// // ShowError displays an error message to the user.
	// ShowError(message string) error
}
