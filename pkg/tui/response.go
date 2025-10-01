package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/SimiyuWafulah/apiterm/internal"
)

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