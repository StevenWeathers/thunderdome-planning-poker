package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/config"
)

var validateConfig bool

var configCmd = &cobra.Command{
	Use:          "config",
	Short:        "Inspect or validate application configuration",
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	Run:          runConfig,
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVar(&validateConfig, "validate", false, "Validate config values required for secure production use")
}

func runConfig(cmd *cobra.Command, args []string) {
	if !validateConfig {
		_ = cmd.Help()
		return
	}

	version := RootCmd.Version
	zlog, _ := zap.NewProduction(
		zap.Fields(
			zap.String("version", version),
		),
	)
	defer func() {
		_ = zlog.Sync()
	}()

	logger := otelzap.New(zlog)
	c := config.InitConfig(logger)
	issues := c.Validate()

	if len(issues) == 0 {
		cmd.Println("configuration validation passed")
		return
	}

	cmd.PrintErrln("configuration validation failed:")
	for _, issue := range issues {
		cmd.PrintErrln(fmt.Sprintf("- %s: %s", issue.Key, issue.Message))
	}

	os.Exit(1)
}
