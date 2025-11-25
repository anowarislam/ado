package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anowarislam/ado/cmd/ado/echo"
	"github.com/anowarislam/ado/cmd/ado/meta"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().String("config", "", "Path to config file")
	cmd.PersistentFlags().String("log-level", "info", "Log level for output")

	cmd.AddCommand(
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
