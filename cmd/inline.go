package cmd

import (
	"fmt"

	"github.com/igoreulalio/container-compliance-checker/internal/config"
	"github.com/igoreulalio/container-compliance-checker/internal/logger"
	"github.com/igoreulalio/container-compliance-checker/internal/service/inline"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// inlineCmd represents the command to run checks inside the container
var inlineCmd = &cobra.Command{
	Use:   "inline",
	Short: "Run compliance checks inside the container",
	Long:  `Executes all configured compliance checks within the running container environment.`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load configuration")
		}

		level, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("Failed to parse log level: %s", err))
		}
		log.Debug().Msgf("Utilizing log level: %s", level)
		logger.InitGlobalLogger(level)

		log.Info().Msg("Starting compliance checks in inline mode")
		inlineService := inline.NewInline(*cfg)
		err = inlineService.Run()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run inline compliance checks")
		}
		log.Info().Msg("Inline compliance checks completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(inlineCmd)
}
