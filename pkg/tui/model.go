package tui

import (
	"github.com/SimiyuWafulah/apiterm/internal/formatter"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

// Model manages the state of our TUI
type Model struct {
	urlInput    textinput.Model
	methodInput textinput.Model
	bodyInput   textinput.Model
	focusIndex  int
	response    string
	loading     bool
	formatter   *formatter.ResponseFormatter

	// terminal size for responsive layout / centering
	width  int
	height int
}

// NewModel creates a new model with default values
func NewModel() *Model {
	m := &Model{
		formatter: formatter.NewResponseFormatter(),
	}

	// URL input
	m.urlInput = textinput.New()
	m.urlInput.Placeholder = "https://apiterm.com/resource"
	m.urlInput.Focus()
	m.urlInput.Width = 60
	m.urlInput.Prompt = " "

	// Method input
	m.methodInput = textinput.New()
	m.methodInput.Placeholder = "GET"
	m.methodInput.CharLimit = 10
	m.methodInput.Width = 10
	m.methodInput.Prompt = " "

	// Body input
	m.bodyInput = textinput.New()
	m.bodyInput.Placeholder = `{"key":"value"}`
	m.bodyInput.Width = 60
	m.bodyInput.Prompt = " "

	// focus the first field
	m.focusIndex = 0
	m.urlInput.Focus()

	// sensible defaults until we get a WindowSizeMsg
	m.width = 80
	m.height = 24

	return m
}

// Init is called when the program starts
func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Run starts application
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}