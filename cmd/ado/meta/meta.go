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
		newSystemCommand(),
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

func newSystemCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "system",
		Short: "Show system diagnostic information",
		Long: `Display system-level diagnostic information including OS, CPU, GPU, NPU, memory, and storage.

Useful for:
  - Troubleshooting environment-specific issues
  - Sharing system information in bug reports
  - Validating system requirements for ado commands
  - Capturing system state in CI/CD pipelines

Output formats:
  - text (default): Human-readable sectioned output
  - json: Structured JSON for parsing/automation
  - yaml: Structured YAML for parsing/automation

Examples:
  # Show system info in human-readable format
  ado meta system

  # Export as JSON for bug report
  ado meta system --output json

  # Extract specific field with jq
  ado meta system --output json | jq '.memory.used_percent'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			info := internalmeta.CollectSystemInfo(ctx)
			format, err := ui.ParseOutputFormat(output)
			if err != nil {
				return err
			}

			return ui.PrintOutput(cmd.OutOrStdout(), format, info, func() (string, error) {
				return formatSystemInfo(info), nil
			})
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json, yaml")
	return cmd
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

func formatSystemInfo(info internalmeta.SystemInfo) string {
	var b strings.Builder

	// OS Section
	fmt.Fprintf(&b, "OS: %s\n", info.OS)
	fmt.Fprintf(&b, "Platform: %s\n", info.Platform)
	fmt.Fprintf(&b, "Kernel: %s\n", info.Kernel)
	fmt.Fprintf(&b, "Architecture: %s\n", info.Architecture)
	fmt.Fprintln(&b)

	// CPU Section
	fmt.Fprintln(&b, "CPU:")
	fmt.Fprintf(&b, "  Model: %s\n", info.CPU.Model)
	fmt.Fprintf(&b, "  Vendor: %s\n", info.CPU.Vendor)
	fmt.Fprintf(&b, "  Cores: %d\n", info.CPU.Cores)
	if info.CPU.FrequencyMHz > 0 {
		fmt.Fprintf(&b, "  Frequency: %.0f MHz\n", info.CPU.FrequencyMHz)
	} else {
		fmt.Fprintln(&b, "  Frequency: unknown")
	}
	fmt.Fprintln(&b)

	// Memory Section
	fmt.Fprintln(&b, "Memory:")
	fmt.Fprintf(&b, "  Total: %d MB\n", info.Memory.TotalMB)
	fmt.Fprintf(&b, "  Available: %d MB\n", info.Memory.AvailableMB)
	fmt.Fprintf(&b, "  Used: %d MB (%.1f%%)\n", info.Memory.UsedMB, info.Memory.UsedPercent)
	if info.Memory.SwapTotalMB > 0 {
		fmt.Fprintf(&b, "  Swap: %d MB total, %d MB used\n", info.Memory.SwapTotalMB, info.Memory.SwapUsedMB)
	}
	fmt.Fprintln(&b)

	// Storage Section
	if len(info.Storage) > 0 {
		fmt.Fprintln(&b, "Storage:")
		for _, storage := range info.Storage {
			fmt.Fprintf(&b, "  %s: %d MB total, %d MB used (%.1f%%)\n",
				storage.Mountpoint, storage.TotalMB, storage.UsedMB, storage.UsedPercent)
		}
		fmt.Fprintln(&b)
	}

	// GPU Section
	if len(info.GPU) > 0 {
		fmt.Fprintln(&b, "GPU:")
		for _, gpu := range info.GPU {
			fmt.Fprintf(&b, "  %s %s (%s)\n", gpu.Vendor, gpu.Model, gpu.Type)
		}
		fmt.Fprintln(&b)
	}

	// NPU Section
	if info.NPU != nil && info.NPU.Detected {
		fmt.Fprintln(&b, "NPU:")
		fmt.Fprintf(&b, "  Type: %s\n", info.NPU.Type)
		fmt.Fprintf(&b, "  Detection Method: %s\n", info.NPU.InferenceMethod)
	}

	return b.String()
}
