package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8")) // Gray
	headerStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			BorderStyle(baseStyle.GetBorderStyle()).
			BorderBottom(true).
			BorderForeground(baseStyle.GetBorderBottomForeground())
	rowStyle = lipgloss.NewStyle().
			Padding(0, 1)
	greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	redStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red
)

func printTable(headers []string, data [][]any) {
	// Convert data to table rows
	rows := make([][]string, len(data))
	for i, row := range data {
		formattedRow := make([]string, len(row))
		for j, cell := range row {
			// Apply cell data formatting based on the type of data
			switch v := cell.(type) {
			case bool:
				formattedRow[j] = formatBool(v)
			case []string:
				formattedRow[j] = formatStringArray(v)
			default:
				formattedRow[j] = fmt.Sprintf("%v", v)
			}
		}
		rows[i] = formattedRow
	}

	// Calculate max width
	columnWidths := make([]int, len(headers))
	for i, header := range headers {
		columnWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}

	totalWidth := 0
	for _, width := range columnWidths {
		totalWidth += width
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			default:
				return rowStyle
			}
		}).
		Headers(headers...).
		Rows(rows...)

	terminalWidth := getTerminalWidth()
	if totalWidth > terminalWidth {
		t.Width(terminalWidth)
	}

	// Print the table
	fmt.Println(t)
}

// Format string array values returning "not present" if empty
func formatStringArray(arr []string) string {
	if len(arr) == 0 {
		return fmt.Sprintf("%s not present", greenStyle.Render("✔"))
	}
	return fmt.Sprintf("%s %s", redStyle.Render("✘"), strings.Join(arr, ", "))
}

// Format boolean values with color where true is bad
func formatBool(value bool) string {
	if value {
		return fmt.Sprintf("%s present", redStyle.Render("✘"))
	}
	return fmt.Sprintf("%s not present", greenStyle.Render("✔"))
}

// getTerminalWidth gets the current terminal width
func getTerminalWidth() int {
	width, _, err := term.GetSize(os.Stdin.Fd())
	if err != nil {
		return 100 // Default fallback
	}
	return width
}
