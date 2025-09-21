package tui

import (
	"fmt"
	"strings"

	"github.com/SimiyuWafulah/apiterm/internal"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbletea/textinput"
	"github.com/charmbracelet/lipgloss"
)

// Model manages the state of our TUI
type Model struct {
	urlInput    textinput.Model
	methodInput textinput.Model
	bodyInput   textinput.Model
	focusIndex  int
	response    string
	err         error
}

// InitialModel creates a new model with default values
func InitialModel() Model {
	m := Model{}

	// URL input
	m.urlInput = textinput.New()
	m.urlInput.Placeholder = "https://apiterm.com/resource"
	m.urlInput.Focus()
	m.urlInput.CharLimit = 200
	m.urlInput.Width = 60

	// Method input
	m.methodInput = textinput.New()
	m.methodInput.Placeholder = "GET"
	m.methodInput.CharLimit = 10
	m.methodInput.Width = 10

	// Body input
	m.bodyInput = textinput.New()
	m.bodyInput.Placeholder = `{"key": "value"}`
	m.bodyInput.Width = 60

	return m
}

// Init is called when the program starts
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			// Handle focus switching
			s := msg.String()
			if s == "enter" && m.focusIndex == 2 {
				// Last field - execute request
				return m, m.executeRequest
			}
			if s == "enter" || s == "down" {
				m.focusIndex++
			} else {
				m.focusIndex--
			}
			if m.focusIndex > 2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 2
			}
			cmds := make([]tea.Cmd, 3)
			cmds[0] = m.urlInput.Blur
			cmds[1] = m.methodInput.Blur
			cmds[2] = m.bodyInput.Blur
			switch m.focusIndex {
			case 0:
				cmds[0] = m.urlInput.Focus
			case 1:
				cmds[1] = m.methodInput.Focus
			case 2:
				cmds[2] = m.bodyInput.Focus
			}
			return m, tea.Batch(cmds...)
		}
	}

	// Update the focused input field
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.urlInput, cmd = m.urlInput.Update(msg)
	case 1:
		m.methodInput, cmd = m.methodInput.Update(msg)
	case 2:
		m.bodyInput, cmd = m.bodyInput.Update(msg)
	}
	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	var b strings.Builder

	// Style for labels
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Width(12).Bold(true)

	// Build the form
	b.WriteString("APITERM - API Client\n\n")
	
	b.WriteString(labelStyle.Render("URL") + ": ")
	b.WriteString(m.urlInput.View())
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Method") + ": ")
	b.WriteString(m.methodInput.View())
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Body") + ": ")
	b.WriteString(m.bodyInput.View())
	b.WriteString("\n\n")

	// Shows response if we have one
	if m.response != "" {
		b.WriteString("Response:\n")
		b.WriteString(m.response)
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n\n")
	b.WriteString("(tab/shift+tab to navigate, enter to send, q to quit)")

	return b.String()
}

// executeRequest runs the HTTP request
func (m Model) executeRequest() tea.Msg {
	url := m.urlInput.Value()
	method := strings.ToUpper(m.methodInput.Value())
	body := m.bodyInput.Value()

	var status int
	var responseBody []byte
	var err error

	switch method {
	case "GET":
		status, responseBody, err = internal.Get(url)
	case "POST":
		status, responseBody, err = internal.Post(url, []byte(body))
	default:
		return fmt.Sprintf("Error: Unsupported method %s", method)
	}

	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return fmt.Sprintf("Status: %d\nBody: %s", status, string(responseBody))
}

// Run starts application
func Run() error {
	p := tea.NewProgram(InitialModel())
	_, err := p.Run()
	return err
}
