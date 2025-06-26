package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Test that root command is properly initialized
	if rootCmd.Use != "elasticdump" {
		t.Errorf("Expected root command use to be 'elasticdump', got '%s'", rootCmd.Use)
	}

	if rootCmd.Version != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got '%s'", rootCmd.Version)
	}

	// Test that subcommands are registered
	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	expectedCommands := []string{"transfer", "backup", "restore"}
	for _, cmdName := range expectedCommands {
		if !commandNames[cmdName] {
			t.Errorf("Expected command '%s' to be registered", cmdName)
		}
	}
}

func TestTransferCommand(t *testing.T) {
	// Test that transfer command is properly configured
	if transferCmd.Use != "transfer" {
		t.Errorf("Expected transfer command use to be 'transfer', got '%s'", transferCmd.Use)
	}

	// Test required flags
	requiredFlags := []string{"input", "output"}
	for _, flagName := range requiredFlags {
		annotations := transferCmd.Annotations
		if annotations == nil {
			// Check if flag is marked as required in another way
			flag := transferCmd.Flag(flagName)
			if flag == nil {
				t.Errorf("Expected flag '%s' to exist on transfer command", flagName)
			}
		}
	}
}

func TestBackupCommand(t *testing.T) {
	// Test that backup command is properly configured
	if backupCmd.Use != "backup" {
		t.Errorf("Expected backup command use to be 'backup', got '%s'", backupCmd.Use)
	}

	// Test that it has input and output flags
	inputFlag := backupCmd.Flag("input")
	if inputFlag == nil {
		t.Errorf("Expected backup command to have 'input' flag")
	}

	outputFlag := backupCmd.Flag("output")
	if outputFlag == nil {
		t.Errorf("Expected backup command to have 'output' flag")
	}
}

func TestRestoreCommand(t *testing.T) {
	// Test that restore command is properly configured
	if restoreCmd.Use != "restore" {
		t.Errorf("Expected restore command use to be 'restore', got '%s'", restoreCmd.Use)
	}

	// Test that it has input and output flags
	inputFlag := restoreCmd.Flag("input")
	if inputFlag == nil {
		t.Errorf("Expected restore command to have 'input' flag")
	}

	outputFlag := restoreCmd.Flag("output")
	if outputFlag == nil {
		t.Errorf("Expected restore command to have 'output' flag")
	}
}

func TestCommandHierarchy(t *testing.T) {
	// Test that all commands are properly added to root
	rootCommands := rootCmd.Commands()

	found := make(map[string]*cobra.Command)
	for _, cmd := range rootCommands {
		found[cmd.Name()] = cmd
	}

	if _, exists := found["transfer"]; !exists {
		t.Error("transfer command not found in root command")
	}

	if _, exists := found["backup"]; !exists {
		t.Error("backup command not found in root command")
	}

	if _, exists := found["restore"]; !exists {
		t.Error("restore command not found in root command")
	}
}
