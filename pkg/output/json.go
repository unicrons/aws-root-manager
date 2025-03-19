package output

import (
	"encoding/json"
	"fmt"
)

// prints data in JSON format
func PrintJSON(headers []string, data [][]any) error {
	// Convert [][]string to a slice of maps for JSON output
	var jsonData []map[string]any
	for _, row := range data {
		rowMap := make(map[string]any)
		for i, cell := range row {
			rowMap[headers[i]] = cell
		}
		jsonData = append(jsonData, rowMap)
	}

	encodedData, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(encodedData))
	return nil
}
