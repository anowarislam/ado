package meta

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	internalmeta "github.com/anowarislam/ado/internal/meta"
	"github.com/anowarislam/ado/internal/ui"
)

func NewCommand(buildInfo internalmeta.BuildInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta",
		Short: "Introspect the ado binary and its environment",
	}

	cmd.AddCommand(
		newInfoCommand(buildInfo),
		newEnvCommand(),
		newFeaturesCommand(),
	)

	return cmd
}

func newInfoCommand(buildInfo internalmeta.BuildInfo) *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show ado build metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := ui.ParseOutputFormat(output)
			if err != nil {
				return err
			}

			return ui.PrintOutput(cmd.OutOrStdout(), format, buildInfo, func() (string, error) {
				return formatBuildInfo(buildInfo), nil
			})
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json, yaml")
	return cmd
}

func newEnvCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Show configuration and environment information",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Root().PersistentFlags().GetString("config")
			if err != nil {
				return err
			}

			info := internalmeta.CollectEnvInfo(configPath)
			format, err := ui.ParseOutputFormat(output)
			if err != nil {
				return err
			}

			return ui.PrintOutput(cmd.OutOrStdout(), format, info, func() (string, error) {
				return formatEnvInfo(info), nil
			})
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json, yaml")
	return cmd
}

func newFeaturesCommand() *cobra.Command {
	var output string
	features := []string{}

	cmd := &cobra.Command{
		Use:   "features",
		Short: "List compiled-in feature flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := ui.ParseOutputFormat(output)
			if err != nil {
				return err
			}

			payload := map[string][]string{"features": features}
			return ui.PrintOutput(cmd.OutOrStdout(), format, payload, func() (string, error) {
				if len(features) == 0 {
					return "No experimental features enabled", nil
				}
				return strings.Join(features, "\n"), nil
			})
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json, yaml")
	return cmd
}

func formatBuildInfo(info internalmeta.BuildInfo) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", info.Name)
	fmt.Fprintf(&b, "Version: %s\n", info.Version)
	fmt.Fprintf(&b, "Commit: %s\n", info.Commit)
	fmt.Fprintf(&b, "BuildTime: %s\n", info.BuildTime)
	fmt.Fprintf(&b, "GoVersion: %s\n", info.GoVersion)
	fmt.Fprintf(&b, "Platform: %s\n", info.Platform)
	return b.String()
}

func formatEnvInfo(info internalmeta.EnvInfo) string {
	var b strings.Builder

	configPath := info.ConfigPath
	if configPath == "" {
		configPath = "(none resolved)"
	}

	fmt.Fprintf(&b, "ConfigPath: %s\n", configPath)
	fmt.Fprintln(&b, "ConfigSources:")
	if len(info.ConfigSources) == 0 {
		fmt.Fprintln(&b, "  (none)")
	} else {
		for _, src := range info.ConfigSources {
			fmt.Fprintf(&b, "  - %s\n", src)
		}
	}

	fmt.Fprintf(&b, "HomeDir: %s\n", info.HomeDir)
	fmt.Fprintf(&b, "CacheDir: %s\n", info.CacheDir)

	fmt.Fprintln(&b, "EnvVariables:")
	if len(info.Env) == 0 {
		fmt.Fprintln(&b, "  (none set)")
	} else {
		for key, value := range info.Env {
			fmt.Fprintf(&b, "  %s=%s\n", key, value)
		}
	}

	return b.String()
}
