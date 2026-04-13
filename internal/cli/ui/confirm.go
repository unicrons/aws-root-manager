package ui

// Confirm shows a yes/no single-select TUI. Returns true if the user chose "Yes".
func Confirm(question string) (bool, error) {
	idx, err := PromptSingle(question, []string{"Yes", "No"})
	if err != nil {
		return false, err
	}
	return idx == 0, nil
}
