package output

import (
	"fmt"
	"io"
	"log/slog"
)

// HandleOutput handles the output based on the specified format
func HandleOutput(w io.Writer, format string, headers []string, rawData [][]any) {
	switch format {
	case "json":
		if err := PrintJSON(w, headers, rawData); err != nil {
			slog.Error("error printing json", "error", err)
		}
	case "csv":
		if err := printCSV(w, headers, dataToString(rawData)); err != nil {
			slog.Error("error printing csv", "error", err)
		}
	case "table":
		printTable(w, headers, rawData)
	default:
		slog.Error("unsupported output format", "format", format)
	}
}

func dataToString(data [][]any) [][]string {
	convertedData := make([][]string, len(data))
	for i, row := range data {
		convertedRow := make([]string, len(row))
		for j, cell := range row {
			convertedRow[j] = fmt.Sprintf("%v", cell)
		}
		convertedData[i] = convertedRow
	}
	return convertedData
}
