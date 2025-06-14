package report

import "time"

type Reporter struct {
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

func NewReporter(results []*Result) *Reporter {

	success := true
	for _, result := range results {
		if result.pass != true {
			success = false
		}
	}

	return &Reporter{
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
