package echo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anowarislam/ado/internal/ui"
)

func NewCommand() *cobra.Command {
	var (
		upper  bool
		lower  bool
		repeat int
		output string
	)

	cmd := &cobra.Command{
		Use:   "echo [message...]",
		Short: "Echo input text with optional formatting",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if upper && lower {
				return errors.New("cannot use --upper and --lower together")
			}
			if repeat < 1 {
				return fmt.Errorf("--repeat must be >= 1 (got %d)", repeat)
			}

			format, err := ui.ParseOutputFormat(output)
			if err != nil {
				return err
			}

			message := strings.Join(args, " ")
			if upper {
				message = strings.ToUpper(message)
			}
			if lower {
				message = strings.ToLower(message)
			}

			values := make([]string, repeat)
			for i := 0; i < repeat; i++ {
				values[i] = message
			}

			return ui.PrintOutput(cmd.OutOrStdout(), format, values, func() (string, error) {
				return strings.Join(values, "\n"), nil
			})
		},
	}

	cmd.Flags().BoolVar(&upper, "upper", false, "Convert message to uppercase")
	cmd.Flags().BoolVar(&lower, "lower", false, "Convert message to lowercase")
	cmd.Flags().IntVar(&repeat, "repeat", 1, "Number of times to repeat the message")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json, yaml")

	return cmd
}
