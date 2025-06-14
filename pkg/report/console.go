package report

import (
	"fmt"
	"strings"
	"time"
)

// PrintConsole displays the report results in a tabular format on the console
func (r *Reporter) PrintConsole() {
	// Print header
	fmt.Println("\n=== Container Compliance Check Reporter ===")
	fmt.Printf("Time: %s\n", r.Timestamp.Format(time.RFC1123))
	fmt.Printf("Overall Status: %s\n\n", formatStatus(r.Success))

	// Define column widths
	checkNameWidth := 40
	passWidth := 10
	errorWidth := 10

	// Print table header
	headerFormat := "%-*s %-*s %-*s %s\n"
	fmt.Printf(headerFormat,
		checkNameWidth, "Check Name",
		passWidth, "Status",
		errorWidth, "Error",
		"Message")

	// Print separator line
	separatorLine := strings.Repeat("-", checkNameWidth+passWidth+errorWidth+40)
	fmt.Println(separatorLine)

	// Print each result
	rowFormat := "%-*s %-*s %-*s %s\n"
	for _, result := range r.Results {
		fmt.Printf(rowFormat,
			checkNameWidth, truncate(result.checkName, checkNameWidth-2),
			passWidth, formatStatus(result.pass),
			errorWidth, formatBoolean(result.haveError),
			truncateWithEllipsis(result.errorString, 60))
	}

	// Print summary
	fmt.Println(separatorLine)
	totalChecks := len(r.Results)
	passedChecks := countPassedChecks(r.Results)
	fmt.Printf("\nSummary: %d/%d checks passed (%d failed)\n\n",
		passedChecks, totalChecks, totalChecks-passedChecks)
}

// formatStatus returns a formatted string for pass/fail status
func formatStatus(pass bool) string {
	if pass {
		return "PASS"
	}
	return "FAIL"
}

// formatBoolean returns a formatted string for boolean values
func formatBoolean(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

// countPassedChecks counts the number of passed checks
func countPassedChecks(results []*Result) int {
	count := 0
	for _, result := range results {
		if result.pass {
			count++
		}
	}
	return count
}

// truncate ensures a string doesn't exceed the maximum length
func truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

// truncateWithEllipsis truncates a string and adds ellipsis if it exceeds the maximum length
func truncateWithEllipsis(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}
