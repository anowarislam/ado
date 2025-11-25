package root

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()

	if cmd.Use != "ado" {
		t.Errorf("Use = %q, want %q", cmd.Use, "ado")
	}

	// Verify subcommands are registered
	subcommands := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = true
	}

	expectedCmds := []string{"echo", "meta"}
	for _, name := range expectedCmds {
		if !subcommands[name] {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}

func TestRootCommand_Help(t *testing.T) {
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ado") {
		t.Errorf("help output missing 'ado'")
	}
	if !strings.Contains(output, "echo") {
		t.Errorf("help output missing 'echo' command")
	}
	if !strings.Contains(output, "meta") {
		t.Errorf("help output missing 'meta' command")
	}
}

func TestRootCommand_Version(t *testing.T) {
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ado") {
		t.Errorf("version output = %q, expected to contain 'ado'", output)
	}
}

func TestRootCommand_GlobalFlags(t *testing.T) {
	cmd := NewRootCommand()

	// Check --config flag exists
	configFlag := cmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("--config flag not found")
	}

	// Check --log-level flag exists
	logLevelFlag := cmd.PersistentFlags().Lookup("log-level")
	if logLevelFlag == nil {
		t.Error("--log-level flag not found")
	}
}
