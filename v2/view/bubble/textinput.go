package bubble

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/schwarzit/go-template/v2/gotemplate"
)

type (
	errMsg error
)

type textInputModel struct {
	textInput textinput.Model
	err       error
}

func initialModel(question gotemplate.TemplateQuestion) textInputModel {
	ti := textinput.New()
	if question.DefaultValue != nil {
		ti.Placeholder = *question.DefaultValue
	}
	ti.Focus()
	// ti.CharLimit = 156
	// ti.Width = 20

	return textInputModel{
		textInput: ti,
		err:       nil,
	}
}

func (m textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	return m.textInput.View()
}
