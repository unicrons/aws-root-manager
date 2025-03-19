package output

import (
	"fmt"

	"github.com/unicrons/aws-root-manager/pkg/logger"
)

// HandleOutput handles the output based on the specified format
func HandleOutput(format string, headers []string, rawData [][]any) {
	switch format {
	case "json":
		if err := PrintJSON(headers, rawData); err != nil {
			logger.Error("output.HandleOutput", err, "error printing json")
		}
	case "csv":
		if err := printCSV(headers, dataToString(rawData)); err != nil {
			logger.Error("output.HandleOutput", err, "error printing csv")
		}
	case "table":
		printTable(headers, rawData)
	default:
		logger.Error("output.HandleOutput", nil, "unsupported output format: %v", format)
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
