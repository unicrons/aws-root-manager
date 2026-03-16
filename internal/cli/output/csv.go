package output

import (
	"encoding/csv"
	"io"
)

// Prints data in CSV format
func printCSV(w io.Writer, headers []string, data [][]string) error {
	writer := csv.NewWriter(w)
	if err := writer.Write(headers); err != nil {
		return err
	}
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
