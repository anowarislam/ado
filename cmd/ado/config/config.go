package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	internalconfig "github.com/anowarislam/ado/internal/config"
	"github.com/anowarislam/ado/internal/ui"
)

// NewCommand returns the config parent command with subcommands.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage ado configuration",
	}

	cmd.AddCommand(
		newValidateCommand(),
	)

	return cmd
}

func newValidateCommand() *cobra.Command {
	var (
		filePath string
		strict   bool
		output   string
	)

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration file",
		Long:  "Validate a configuration file against the expected schema and report errors.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Resolve config path
			path := filePath
			if path == "" {
				// Try --config flag from root
				configFlag, _ := cmd.Root().PersistentFlags().GetString("config")
				if configFlag != "" {
					path = configFlag
				}
			}

			if path == "" {
				// Auto-detect
				homeDir, _ := os.UserHomeDir()
				resolved, sources := internalconfig.ResolveConfigPath("", homeDir)
				if resolved == "" {
					return fmt.Errorf("no config file found. Searched: %s", strings.Join(sources, ", "))
				}
				path = resolved
			}

			// Validate
			result, err := internalconfig.Validate(path)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// In strict mode, warnings become errors
			if strict && result.HasWarnings() {
				for _, w := range result.Warnings {
					result.Errors = append(result.Errors, internalconfig.ValidationIssue{
						Message:  w.Message,
						Line:     w.Line,
						Severity: "error",
					})
				}
				result.Warnings = []internalconfig.ValidationIssue{}
				result.Valid = false
			}

			// Output
			format, err := ui.ParseOutputFormat(output)
			if err != nil {
				return err
			}

			err = ui.PrintOutput(cmd.OutOrStdout(), format, result, func() (string, error) {
				return formatValidationResult(result), nil
			})
			if err != nil {
				return err
			}

			// Exit with error code if invalid
			if !result.Valid {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to config file to validate")
	cmd.Flags().BoolVarP(&strict, "strict", "s", false, "Treat warnings as errors")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json")

	return cmd
}

func formatValidationResult(result *internalconfig.ValidationResult) string {
	var b strings.Builder

	if result.Valid {
		fmt.Fprintf(&b, "\u2713 Config valid: %s", result.Path)
	} else {
		fmt.Fprintf(&b, "\u2717 Config invalid: %s", result.Path)
	}

	for _, e := range result.Errors {
		b.WriteString("\n")
		if e.Line > 0 {
			fmt.Fprintf(&b, "  Error: %s at line %d", e.Message, e.Line)
		} else {
			fmt.Fprintf(&b, "  Error: %s", e.Message)
		}
	}

	for _, w := range result.Warnings {
		b.WriteString("\n")
		if w.Line > 0 {
			fmt.Fprintf(&b, "  Warning: %s at line %d", w.Message, w.Line)
		} else {
			fmt.Fprintf(&b, "  Warning: %s", w.Message)
		}
	}

	return b.String()
}
