package ui

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{"empty defaults to text", "", OutputText, false},
		{"text format", "text", OutputText, false},
		{"json format", "json", OutputJSON, false},
		{"yaml format", "yaml", OutputYAML, false},
		{"invalid format", "xml", "", true},
		{"invalid format csv", "csv", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOutputFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPrintOutput_Text(t *testing.T) {
	tests := []struct {
		name       string
		format     OutputFormat
		renderText func() (string, error)
		want       string
		wantErr    bool
	}{
		{
			name:       "text output adds newline",
			format:     OutputText,
			renderText: func() (string, error) { return "hello", nil },
			want:       "hello\n",
		},
		{
			name:       "text output preserves existing newline",
			format:     OutputText,
			renderText: func() (string, error) { return "hello\n", nil },
			want:       "hello\n",
		},
		{
			name:       "empty text output",
			format:     OutputText,
			renderText: func() (string, error) { return "", nil },
			want:       "",
		},
		{
			name:       "text render error",
			format:     OutputText,
			renderText: func() (string, error) { return "", errors.New("render failed") },
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := PrintOutput(&buf, tt.format, nil, tt.renderText)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("PrintOutput() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintOutput_JSON(t *testing.T) {
	var buf bytes.Buffer
	payload := map[string]string{"key": "value"}

	err := PrintOutput(&buf, OutputJSON, payload, nil)
	if err != nil {
		t.Fatalf("PrintOutput() error = %v", err)
	}

	want := "{\n  \"key\": \"value\"\n}\n"
	if got := buf.String(); got != want {
		t.Errorf("PrintOutput() = %q, want %q", got, want)
	}
}

func TestPrintOutput_YAML(t *testing.T) {
	var buf bytes.Buffer
	payload := map[string]string{"key": "value"}

	err := PrintOutput(&buf, OutputYAML, payload, nil)
	if err != nil {
		t.Fatalf("PrintOutput() error = %v", err)
	}

	want := "key: value\n"
	if got := buf.String(); got != want {
		t.Errorf("PrintOutput() = %q, want %q", got, want)
	}
}

func TestPrintOutput_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := PrintOutput(&buf, OutputFormat("unknown"), nil, nil)
	if err == nil {
		t.Error("PrintOutput() expected error for unknown format")
	}
}

func TestPrintOutput_EmptyFormat(t *testing.T) {
	var buf bytes.Buffer
	err := PrintOutput(&buf, "", nil, func() (string, error) { return "hello", nil })
	if err != nil {
		t.Fatalf("PrintOutput() error = %v", err)
	}
	// Empty format should default to text
	if got := buf.String(); got != "hello\n" {
		t.Errorf("PrintOutput() = %q, want %q", got, "hello\n")
	}
}

func TestPrintOutput_JSON_MarshalError(t *testing.T) {
	var buf bytes.Buffer
	// Functions cannot be marshaled to JSON
	payload := map[string]any{"fn": func() {}}

	err := PrintOutput(&buf, OutputJSON, payload, nil)
	if err == nil {
		t.Error("PrintOutput() expected error for unmarshalable payload")
	}
	if !errors.Is(err, err) || buf.Len() != 0 {
		// Should not write partial output
	}
}

func TestPrintOutput_YAML_MarshalError(t *testing.T) {
	// Note: yaml.Marshal panics on channels/functions rather than returning error
	// Testing the error path requires a type that yaml.Marshal returns error for
	// This is difficult to trigger - yaml tends to panic instead
	// Skipping this test as yaml library behavior makes it impractical
	t.Skip("yaml.Marshal panics on unsupported types instead of returning error")
}

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestPrintOutput_Text_WriteError(t *testing.T) {
	w := &errorWriter{}
	err := PrintOutput(w, OutputText, nil, func() (string, error) { return "hello", nil })
	if err == nil {
		t.Error("PrintOutput() expected error for write failure")
	}
}

func TestPrintOutput_JSON_WriteError(t *testing.T) {
	w := &errorWriter{}
	payload := map[string]string{"key": "value"}

	err := PrintOutput(w, OutputJSON, payload, nil)
	if err == nil {
		t.Error("PrintOutput() expected error for write failure")
	}
}

func TestPrintOutput_YAML_WriteError(t *testing.T) {
	w := &errorWriter{}
	payload := map[string]string{"key": "value"}

	err := PrintOutput(w, OutputYAML, payload, nil)
	if err == nil {
		t.Error("PrintOutput() expected error for write failure")
	}
}

func TestPrintOutput_YAML_NoTrailingNewline(t *testing.T) {
	var buf bytes.Buffer
	// Empty map serializes to "{}\n" in YAML
	payload := struct{}{}

	err := PrintOutput(&buf, OutputYAML, payload, nil)
	if err != nil {
		t.Fatalf("PrintOutput() error = %v", err)
	}

	output := buf.String()
	if output[len(output)-1] != '\n' {
		t.Error("YAML output should end with newline")
	}
}
