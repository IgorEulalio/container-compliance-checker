package checksmanager

import (
	"fmt"

	"github.com/igoreulalio/container-compliance-checker/internal/checks"
	"github.com/igoreulalio/container-compliance-checker/internal/checks/file"
	packagecheck "github.com/igoreulalio/container-compliance-checker/internal/checks/package"
	"github.com/igoreulalio/container-compliance-checker/internal/config"
	"github.com/rs/zerolog/log"
)

func MapConfigToChecks(config config.Config) ([]checks.Check, error) {
	var checkList []checks.Check
	var errList []error
	for _, c := range config.Checks {
		switch c.Type {
		case PackageNotInstalled:
			check, err := packagecheck.NewPackageNotInstalledCheck(c.Config)
			if err != nil {
				log.Error().Err(err).Str("check", c.Type).Msg("Error creating check.")
				errList = append(errList, err)
				continue
			}
			checkList = append(checkList, check)
		case FilePermission:
			check, err := file.NewFilePermissionCheck(c.Config)
			if err != nil {
				log.Error().Err(err).Str("check", c.Type).Msg("Error creating check.")
				errList = append(errList, err)
				continue
			}
			checkList = append(checkList, check)
		default:
			log.Warn().Msgf("Unknown check type, skipping: %s", c.Type)
			continue
		}
	}

	if len(errList) > 0 {
		return nil, fmt.Errorf("error creating checks")
	}

	return checkList, nil
}
