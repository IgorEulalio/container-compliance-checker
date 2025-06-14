package report

import (
	"encoding/csv"
	"os"
)

// WriteCSV writes the report results to a CSV file with columns: CheckName, Pass, HaveError, ErrorString
func (r *Report) WriteCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"CheckName", "Pass", "HaveError", "ErrorString"})
	if err != nil {
		return err
	}

	for _, result := range r.Results {
		record := []string{
			result.checkName,
			boolToString(result.pass),
			boolToString(result.haveError),
			result.errorString,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
