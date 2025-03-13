package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func DisplayTable(headers []string, data [][]string) {

	// Convert data to table rows
	rows := make([]table.Row, len(data))
	for i, row := range data {
		rows[i] = row
	}

	// Calculate max width for each column
	columnWidths := make([]int, len(headers))

	// Check headers length
	for i, header := range headers {
		columnWidths[i] = len(header)
	}

	// Check data length for each column
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}

	// Create columns with calculated widths
	columns := make([]table.Column, len(headers))
	for i, header := range headers {
		// Add some padding to the width
		columns[i] = table.Column{
			Title: header,
			Width: columnWidths[i],
		}
	}

	// Customize the table styles
	style := table.Styles{
		Header: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			BorderStyle(baseStyle.GetBorderStyle()).
			BorderBottom(true).
			BorderForeground(baseStyle.GetBorderBottomForeground()),
		Cell: lipgloss.NewStyle().Padding(0, 1),
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
