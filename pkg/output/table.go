package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	redStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red
)

func printTable(headers []string, data [][]any) {
	termWidth := getTerminalWidth()
	maxColWidth := termWidth / len(headers)

	// Convert data to table rows
	rows := make([]table.Row, len(data))
	for i, row := range data {
		formattedRow := make([]string, len(row))
		for j, cell := range row {
			// Apply formatting based on the type of the cell
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

	// Calculate max width for each column
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

	// Create table columns with the adjusted widths
	columns := make([]table.Column, len(headers))
	for i, header := range headers {
		columns[i] = table.Column{
			Title: header,
			Width: min(columnWidths[i], maxColWidth),
		}
	}

	// Customize the table styles
	style := table.Styles{
		Header: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			BorderStyle(baseStyle.GetBorderStyle()).
			BorderBottom(true).
			BorderForeground(baseStyle.GetBorderBottomForeground()).
			MaxWidth(maxColWidth),
		Cell: lipgloss.NewStyle().
			Padding(0, 1).
			MaxWidth(maxColWidth),
	}

	// Create the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)+1),
		table.WithStyles(style),
	)

	// Render the table
	output := baseStyle.Render(t.View())
	fmt.Println(output)
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
