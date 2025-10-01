package tui

import (
	"fmt"
	"strings"

	"github.com/SimiyuWafulah/apiterm/internal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// Model manages the state of our TUI
type Model struct {
	urlInput    textinput.Model
	methodInput textinput.Model
	bodyInput   textinput.Model
	focusIndex  int
	response    string
	loading     bool

	// terminal size for responsive layout / centering
	width  int
	height int
}

// NewModel creates a new model with default values
func NewModel() *Model {
	m := &Model{}

	// URL input
	m.urlInput = textinput.New()
	m.urlInput.Placeholder = "https://apiterm.com/resource"
	m.urlInput.Focus()
	//m.urlInput.CharLimit = 200
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
	//m.bodyInput.CharLimit = 1000
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

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n": 
			if !m.loading {
				m.resetForm()
				return m, nil
			}
		case "tab", "down":
			m.focusIndex = (m.focusIndex + 1) % 3
			m.updateFocus()
			return m, nil
		case "shift+tab", "up":
			m.focusIndex = (m.focusIndex - 1 + 3) % 3
			m.updateFocus()
			return m, nil
		case "enter":
			// if "enter" pressed on last field, send request
			if m.focusIndex == 2 && !m.loading {
				m.loading = true
				return m, tea.Cmd(m.executeRequest)
			}
			// otherwise move focus forward
			m.focusIndex = (m.focusIndex + 1) % 3
			m.updateFocus()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case string:
		// response message
		m.response = msg
		m.loading = false
	}

	// Update the inputs and collect cmds
	m.urlInput, cmd = m.urlInput.Update(msg)
	cmds = append(cmds, cmd)
	m.methodInput, cmd = m.methodInput.Update(msg)
	cmds = append(cmds, cmd)
	m.bodyInput, cmd = m.bodyInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateFocus updates which input has focus based on focusIndex
func (m *Model) updateFocus() {
	switch m.focusIndex {
	case 0:
		m.urlInput.Focus()
		m.methodInput.Blur()
		m.bodyInput.Blur()
	case 1:
		m.urlInput.Blur()
		m.methodInput.Focus()
		m.bodyInput.Blur()
	case 2:
		m.urlInput.Blur()
		m.methodInput.Blur()
		m.bodyInput.Focus()
	}
}

// NEW: resetForm clears all fields for a new request
func (m *Model) resetForm() {
	m.urlInput.SetValue("")
	m.methodInput.SetValue("")
	m.bodyInput.SetValue("")
	m.response = ""
	m.focusIndex = 0
	m.urlInput.Focus()
}

// View renders the ui
func (m *Model) View() string {
	var b strings.Builder

	termWidth := m.width
	if termWidth <= 0 {
		termWidth = 80
	}

	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("205")).
		Width(termWidth).
		Align(lipgloss.Center)

	header := titleStyle.Render("APITERM - API Client")

	// Label and cell styles
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Width(12).Bold(true)
	cellStyle := lipgloss.NewStyle().PaddingLeft(1)

	
	formBoxWidth := termWidth
	if formBoxWidth > 4 {
		formBoxWidth = formBoxWidth - 4
	}
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Width(formBoxWidth)

	
	var rows []string

	row := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("URL")+":",
		cellStyle.Render(m.urlInput.View()),
	)
	rows = append(rows, row)

	row = lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Method")+":",
		cellStyle.Render(m.methodInput.View()),
	)
	rows = append(rows, row)

	row = lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Body")+":",
		cellStyle.Render(m.bodyInput.View()),
	)
	rows = append(rows, row)

	form := strings.Join(rows, "\n\n")

	// Response area styling
	respStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(formBoxWidth)

	
thelpStyle := lipgloss.NewStyle().Italic(true).PaddingTop(1)

	b.WriteString(header + "\n\n")
	b.WriteString(formStyle.Render(form))
	b.WriteString("\n\n")

	if m.response != "" {
		b.WriteString(respStyle.Render("Response:\n" + m.response))
	} else if m.loading {
		b.WriteString(respStyle.Render("Response:\nLoading..."))
	} else {
		b.WriteString(respStyle.Render("Response: (no response yet)"))
	}

	// Updated help text to include 'ctrl+n' for new request
	b.WriteString("\n\n" + thelpStyle.Render("(ctrl+n: new request, tab/shift+tab or up/down to navigate, enter to send on Body, q to quit)"))

	return b.String()
}

// executeRequest runs the HTTP request 
func (m *Model) executeRequest() tea.Msg {
	url := strings.TrimSpace(m.urlInput.Value())
	method := strings.ToUpper(strings.TrimSpace(m.methodInput.Value()))
	body := m.bodyInput.Value()

	if url == "" {
		return "Error: URL is required"
	}
	if method == "" {
		method = "GET"
	}

	var status int
	var responseBody []byte
	var err error

	switch method {
	case "GET":
		status, responseBody, err = internal.Get(url)
	case "POST":
		status, responseBody, err = internal.Post(url, []byte(body))
	default:
		return "Error: Unsupported method " + method
	}

	if err != nil {
		return "Error: " + err.Error()
	}

	return fmt.Sprintf("Status: %d\nBody: %s", status, string(responseBody))
}

// Run starts application
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}