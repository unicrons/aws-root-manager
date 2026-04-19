package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
)

// Confirm shows a yes/no single-select TUI. Returns true if the user chose "Yes".
// Returns an error if stdin is not a TTY — callers must pass --yes to skip confirmation in non-interactive environments.
func Confirm(question string) (bool, error) {
	if !term.IsTerminal(os.Stdin.Fd()) {
		return false, fmt.Errorf("confirmation required: run in an interactive terminal or use --yes to confirm")
	}
	idx, err := PromptSingle(question, []string{"Yes", "No"})
	if err != nil {
		return false, err
	}
	return idx == 0, nil
}
