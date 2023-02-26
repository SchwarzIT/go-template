package bubble

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/schwarzit/go-template/v2/gotemplate"
)

var (
	ErrInterrupt = errors.New("interrupted")
)

// Model represents the state of the view at any given time.
type Model struct {
	// Question is the current question being presented to the user.
	Question gotemplate.TemplateQuestion
	// Input is the user's response to the current question.
	Input string
	// Error is any error that occurred while processing the user's response to the current question.
	Error error
	// Choices is the list of pre-defined choices for the current question, if applicable.
	Choices []string
	// SelectedIdx is the index of the selected choice, if the current question is a multiple choice question.
	SelectedIdx int
}

// BubbleTeaView is a struct that implements the View interface using the Bubble Tea framework.
type BubbleTeaView struct {
	model    Model
	program  *tea.Program
	done     chan struct{}
	response chan string
}

// NewBubbleTeaView creates a new BubbleTeaView.
func NewBubbleTeaView() BubbleTeaView {
	return BubbleTeaView{
		model:    Model{},
		program:  nil,
		done:     make(chan struct{}),
		response: make(chan string),
	}
}

func (b BubbleTeaView) PresentQuestion(question gotemplate.TemplateQuestion) (*gotemplate.TemplateQuestion, error) {
	fmt.Println("PresentQuestion")
	fmt.Println(question)
	question.ResponseValue = nil

	p := tea.NewProgram(initialModel(question))
	if _, err := p.Run(); err != nil {
		return nil, err
	}

	return &question, nil
}

// ShowMessage displays a message to the user.
func (b *BubbleTeaView) ShowMessage(message string) error {
	fmt.Println(message)
	return nil
}
