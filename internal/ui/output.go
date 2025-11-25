package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type OutputFormat string

const (
	OutputText OutputFormat = "text"
	OutputJSON OutputFormat = "json"
	OutputYAML OutputFormat = "yaml"
)

func ParseOutputFormat(raw string) (OutputFormat, error) {
	if raw == "" {
		return OutputText, nil
	}

	switch OutputFormat(raw) {
	case OutputText, OutputJSON, OutputYAML:
		return OutputFormat(raw), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", raw)
	}
}

func PrintOutput(w io.Writer, format OutputFormat, payload any, renderText func() (string, error)) error {
	switch format {
	case OutputText, "":
		text, err := renderText()
		if err != nil {
			return err
		}
		if text == "" {
			return nil
		}
		if text[len(text)-1] != '\n' {
			text += "\n"
		}
		_, err = io.WriteString(w, text)
		return err
	case OutputJSON:
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return fmt.Errorf("serialize json: %w", err)
		}
		_, err = w.Write(append(data, '\n'))
		return err
	case OutputYAML:
		data, err := yaml.Marshal(payload)
		if err != nil {
			return fmt.Errorf("serialize yaml: %w", err)
		}
		_, err = w.Write(data)
		if err != nil {
			return err
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			_, err = w.Write([]byte("\n"))
		}
		return err
	default:
		return errors.New("unknown output format")
	}
}
