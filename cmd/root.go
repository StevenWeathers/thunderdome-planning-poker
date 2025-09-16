package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "thunderdome",
	Short: "Thunderdome is an open source agile tool suite for remote teams.",
	Long: `To get started run the serve subcommand which will start a server
on localhost:8080:

    thunderdome serve
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	RootCmd.Version = version
	if len(os.Args) == 1 {
		// No subcommand given -> run "serve"
		slog.Info("no subcommand given, defaulting to serve")
		os.Args = append(os.Args, "serve")
	}
	if err := RootCmd.Execute(); err != nil {
		slog.Error("error executing root command", slog.Any("error", err))
		os.Exit(-1)
	}
}

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)
}
