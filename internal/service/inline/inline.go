package inline

import (
	"github.com/igoreulalio/container-compliance-checker/internal/config"
	"github.com/igoreulalio/container-compliance-checker/pkg/checksmanager"
	"github.com/igoreulalio/container-compliance-checker/pkg/report"
	"github.com/rs/zerolog/log"
)

type Inline struct {
	Config config.Config
}

// NewInline creates a new Inline service instance with the provided configuration.
func NewInline(cfg config.Config) *Inline {
	return &Inline{
		Config: cfg,
	}
}

// Run executes the compliance checks in inline mode.
func (i *Inline) Run() error {
	checks, err := checksmanager.MapConfigToChecks(i.Config)
	if err != nil {
		return err
	}

	if len(checks) == 0 {
		log.Warn().Msg("No compliance checks configured, skipping execution")
		return nil
	}

	var results []*report.Result
	for _, check := range checks {
		pass, err := check.Run()
		if err != nil {
			results = append(results, report.NewReportResult(check.Name(), false, true, err.Error()))
			continue
		}
		results = append(results, report.NewReportResult(check.Name(), pass, false, ""))
	}
	reporter := report.NewReporter(results)
	reporter.PrintConsole()
	return nil
}
