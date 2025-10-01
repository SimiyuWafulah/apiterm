package tui

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

// resetForm clears all fields for a new request
func (m *Model) resetForm() {
	m.urlInput.SetValue("")
	m.methodInput.SetValue("")
	m.bodyInput.SetValue("")
	m.response = ""
	m.focusIndex = 0
	m.urlInput.Focus()
}