package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

	helpStyle := lipgloss.NewStyle().Italic(true).PaddingTop(1)

	b.WriteString(header + "\n\n")
	b.WriteString(formStyle.Render(form))
	b.WriteString("\n\n")

	if m.response != "" {
		// Use the formatter to format the response
		formattedResponse := m.formatter.FormatResponse(m.response)
		b.WriteString(respStyle.Render("Response:\n" + formattedResponse))
	} else if m.loading {
		b.WriteString(respStyle.Render("Response:\nLoading..."))
	} else {
		b.WriteString(respStyle.Render("Response: (no response yet)"))
	}

	// NEW: Updated help text to include 'ctrl+n' for new request
	b.WriteString("\n\n" + helpStyle.Render("(ctrl+n: new request, tab/shift+tab or up/down to navigate, enter to send on Body, q to quit)"))

	return b.String()
}