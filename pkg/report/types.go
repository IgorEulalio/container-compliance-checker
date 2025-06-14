package report

import "time"

type Report struct {
	Success   bool
	Timestamp time.Time
	Results   []*Result
}

type Result struct {
	checkName   string
	pass        bool
	haveError   bool
	errorString string
}

func NewReport(results []*Result) *Report {

	success := true
	for _, result := range results {
		if result.pass != true {
			success = false
		}
	}

	return &Report{
		Success:   success,
		Timestamp: time.Now(),
		Results:   results,
	}
}

func NewReportResult(checkName string, pass bool, haveError bool, errorString string) *Result {
	return &Result{
		checkName:   checkName,
		pass:        pass,
		haveError:   haveError,
		errorString: errorString,
	}
}
