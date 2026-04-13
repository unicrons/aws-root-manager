package output

import (
	"encoding/json"
	"fmt"
	"io"

	"charm.land/lipgloss/v2"
)

var policyBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("8")).
	Padding(0, 1)

// RenderPolicy pretty-prints a JSON policy string inside a styled box and writes it to w.
// Used for interactive (table) output mode only.
func RenderPolicy(w io.Writer, policy string) {
	pretty, err := prettyJSON(policy)
	if err != nil {
		pretty = policy
	}
	termWidth := getTerminalWidth()
	fmt.Fprintln(w, policyBoxStyle.MaxWidth(termWidth-4).Render(pretty))
}

func prettyJSON(raw string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
