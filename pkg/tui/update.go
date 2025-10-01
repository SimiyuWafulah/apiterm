package tui

import tea "github.com/charmbracelet/bubbletea"

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n": // NEW: Add new request functionality
			if !m.loading {
				m.resetForm()
				// Return early to prevent text inputs from processing the key
				return m, nil
			}
		case "tab", "down":
			m.focusIndex = (m.focusIndex + 1) % 3
			m.updateFocus()
			// Return early to prevent text inputs from processing the key
			return m, nil
		case "shift+tab", "up":
			m.focusIndex = (m.focusIndex - 1 + 3) % 3
			m.updateFocus()
			// Return early to prevent text inputs from processing the key
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
			// Return early to prevent text inputs from processing the key
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