package echo

import (
	"bytes"
	"strings"
	"testing"
)

func TestEchoCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "simple echo",
			args: []string{"hello", "world"},
			want: "hello world\n",
		},
		{
			name: "uppercase",
			args: []string{"--upper", "hello"},
			want: "HELLO\n",
		},
		{
			name: "lowercase",
			args: []string{"--lower", "HELLO"},
			want: "hello\n",
		},
		{
			name: "repeat",
			args: []string{"--repeat", "3", "hi"},
			want: "hi\nhi\nhi\n",
		},
		{
			name:    "upper and lower error",
			args:    []string{"--upper", "--lower", "hello"},
			wantErr: true,
		},
		{
			name:    "no args error",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got := buf.String(); got != tt.want {
					t.Errorf("Execute() output = %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func TestEchoCommand_JSONOutput(t *testing.T) {
	cmd := NewCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--output", "json", "hello"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"hello"`) {
		t.Errorf("Execute() output = %q, want JSON containing 'hello'", got)
	}
}

func TestEchoCommand_YAMLOutput(t *testing.T) {
	cmd := NewCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--output", "yaml", "hello"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "hello") {
		t.Errorf("Execute() output = %q, want YAML containing 'hello'", got)
	}
}
