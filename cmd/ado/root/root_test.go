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

func TestRootCommand_LogLevel_Invalid(t *testing.T) {
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	// Use a subcommand to trigger PersistentPreRunE (--help bypasses it)
	cmd.SetArgs([]string{"--log-level", "invalid", "echo", "test"})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid log level")
		return
	}
	if !strings.Contains(err.Error(), "invalid log level") {
		t.Errorf("error = %q, expected to contain 'invalid log level'", err.Error())
	}
}

func TestRootCommand_LogLevel_Valid(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			cmd := NewRootCommand()
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			// Use echo command to trigger PersistentPreRunE
			cmd.SetArgs([]string{"--log-level", level, "echo", "test"})

			if err := cmd.Execute(); err != nil {
				t.Errorf("Execute() with log level %q error = %v", level, err)
			}
		})
	}
}

func TestRootCommand_ConfigSubcommand(t *testing.T) {
	cmd := NewRootCommand()

	// Verify config subcommand is registered
	subcommands := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = true
	}

	if !subcommands["config"] {
		t.Error("expected subcommand 'config' not found")
	}
}

func TestRootCommand_NoArgs(t *testing.T) {
	// Test running root command with no arguments
	// Should display help (calls cmd.Help() in RunE)
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() with no args error = %v", err)
	}

	output := buf.String()
	// Should show help output
	if !strings.Contains(output, "Usage:") {
		t.Errorf("expected help output with no args, got: %s", output)
	}
	if !strings.Contains(output, "Available Commands:") {
		t.Errorf("expected available commands in help output")
	}
}

func TestRootCommand_PersistentPreRunE(t *testing.T) {
	// Test that PersistentPreRunE initializes logging correctly
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--log-level", "debug", "echo", "test"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify logging was initialized (context should be set)
	// We verify this indirectly by checking the command executed successfully
	if cmd.Context() == nil {
		t.Error("expected context to be set after PersistentPreRunE")
	}
}

func TestRootCommand_UnknownSubcommand(t *testing.T) {
	// Test running unknown subcommand
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for unknown subcommand")
	}

	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("error = %q, expected to contain 'unknown command'", err.Error())
	}
}
