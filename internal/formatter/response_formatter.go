package formatter

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ResponseFormatter handles formatting and syntax highlighting of API responses
type ResponseFormatter struct {
	keyStyle    lipgloss.Style
	stringStyle lipgloss.Style
	numberStyle lipgloss.Style
	boolStyle   lipgloss.Style
	nullStyle   lipgloss.Style
}

// NewResponseFormatter creates a new formatter with default styles
func NewResponseFormatter() *ResponseFormatter {
	return &ResponseFormatter{
		keyStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
		stringStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("85")),
		numberStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("214")),
		boolStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		nullStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true),
	}
}

// FormatResponse takes a raw response string and formats it with syntax highlighting
func (f *ResponseFormatter) FormatResponse(response string) string {
	// Check if the response contains JSON by looking for common patterns
	if strings.Contains(response, "Body: {") || strings.Contains(response, "Body: [") {
		// Extract the body part for JSON formatting
		parts := strings.SplitN(response, "Body: ", 2)
		if len(parts) == 2 {
			body := parts[1]
			
			// Try to parse and pretty-print the JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(body), "", "  "); err == nil {
				// JSON is valid, apply syntax highlighting
				formattedBody := f.applyJSONSyntaxHighlighting(prettyJSON.String())
				return parts[0] + "Body: " + formattedBody
			}
		}
	}
	
	// If not JSON or formatting fails, return original response
	return response
}

// applyJSONSyntaxHighlighting applies colors to JSON components
func (f *ResponseFormatter) applyJSONSyntaxHighlighting(jsonStr string) string {
	result := jsonStr
	lines := strings.Split(result, "\n")
	
	for i, line := range lines {
		// Colorize keys (strings before colon)
		if strings.Contains(line, "\":") {
			parts := strings.SplitN(line, "\":", 2)
			if len(parts) == 2 {
				key := parts[0] + "\":"
				lines[i] = f.keyStyle.Render(key) + parts[1]
			}
		}
		
		// Colorize string values (between quotes, but not keys)
		line = lines[i]
		if strings.Contains(line, ": \"") {
			start := strings.Index(line, ": \"")
			if start != -1 {
				prefix := line[:start+2]
				suffix := line[start+2:]
				
				// Find the closing quote
				if end := strings.Index(suffix, "\""); end != -1 {
					stringVal := suffix[:end+1]
					remaining := suffix[end+1:]
					lines[i] = prefix + f.stringStyle.Render(stringVal) + remaining
				}
			}
		}
		
		// Colorize numbers, booleans, and null
		line = lines[i]
		line = strings.ReplaceAll(line, " true", " "+f.boolStyle.Render("true"))
		line = strings.ReplaceAll(line, " false", " "+f.boolStyle.Render("false"))
		line = strings.ReplaceAll(line, " null", " "+f.nullStyle.Render("null"))
		
		// Simple number detection
		words := strings.Fields(line)
		for _, word := range words {
			if f.isNumber(word) {
				line = strings.ReplaceAll(line, " "+word, " "+f.numberStyle.Render(word))
			}
		}
		lines[i] = line
	}
	
	return strings.Join(lines, "\n")
}

// isNumber checks if a string is a number (basic implementation)
func (f *ResponseFormatter) isNumber(s string) bool {
	if s == "" {
		return false
	}
	// Remove potential trailing commas
	clean := strings.TrimRight(s, ",")
	
	// Check if it's a number (integer or float)
	hasDot := false
	for i, char := range clean {
		if char == '-' && i == 0 {
			continue // Allow negative sign at start
		}
		if char == '.' {
			if hasDot {
				return false // Multiple dots
			}
			hasDot = true
			continue
		}
		if char < '0' || char > '9' {
			return false
		}
	}
	return len(clean) > 0 && (clean != "-" && clean != ".")
}