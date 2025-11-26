package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anowarislam/ado/cmd/ado/config"
	"github.com/anowarislam/ado/cmd/ado/echo"
	"github.com/anowarislam/ado/cmd/ado/meta"
	"github.com/anowarislam/ado/internal/logging"
	internalmeta "github.com/anowarislam/ado/internal/meta"
)

func NewRootCommand() *cobra.Command {
	buildInfo := internalmeta.CurrentBuildInfo()

	cmd := &cobra.Command{
		Use:           "ado",
		Short:         "ado is a composable automation and diagnostics CLI",
		Long:          "ado is a single binary for automation and diagnostics, with discoverable subcommands and consistent UX.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       buildInfo.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Initialize logger from flags
			logLevel, _ := cmd.Flags().GetString("log-level")
			if logLevel != "" && !logging.IsValidLevel(logLevel) {
				return fmt.Errorf("invalid log level %q: must be debug, info, warn, or error", logLevel)
			}

			cfg := logging.Config{
				Level:  logLevel,
				Format: "auto",
				Output: "stderr",
			}.Validate()

			log := logging.New(cfg)
			ctx := logging.WithContext(cmd.Context(), log)
			cmd.SetContext(ctx)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().String("config", "", "Path to config file")
	cmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")

	cmd.AddCommand(
		config.NewCommand(),
		echo.NewCommand(),
		meta.NewCommand(buildInfo),
	)

	return cmd
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
