package output

import (
	"encoding/csv"
	"os"
)

// Prints data in CSV format
func printCSV(headers []string, data [][]string) error {
	writer := csv.NewWriter(os.Stdout)
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
