// Package test_cases provides helpers for shell test assertions.
package test_cases

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/logger"
)

// AssertHistoryFileHasCommands reads the file at filePath, splits it into lines, and checks that each line matches the Command field of the corresponding CommandResponseTestCase.
// It logs a success message for each line and returns a detailed error if any line does not match or if the line count is wrong.
func AssertHistoryFileHasCommands(l *logger.Logger, filePath string, commands []string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n") // TODO: Remove the TrimSpace
	found_commands := []string{}
	for i, command := range commands {
		if i >= len(lines) {
			utils.LogReadableFileContents(l, string(content), fmt.Sprintf("Reading contents from %s", filePath), filePath)
			printFoundCommands(l, found_commands, filePath)
			return fmt.Errorf("file %s has %d lines, expected at least %d", filePath, len(lines), len(commands))
		}
		if lines[i] != command {
			utils.LogReadableFileContents(l, string(content), fmt.Sprintf("Reading contents from %s", filePath), filePath)
			printFoundCommands(l, found_commands, filePath)
			return fmt.Errorf("expected command %q at line %d, got %q", command, i+1, lines[i])
		}
		found_commands = append(found_commands, command)
	}
	if len(lines) != len(commands) {
		utils.LogReadableFileContents(l, string(content), fmt.Sprintf("Reading contents from %s", filePath), filePath)
		printFoundCommands(l, found_commands, filePath)
		return fmt.Errorf("file %s has %d lines, expected %d", filePath, len(lines), len(commands))
	}
	return nil
}

func printFoundCommands(l *logger.Logger, found_commands []string, filePath string) {
	for _, command := range found_commands {
		l.Successf("✓ Found command %q in %s", command, filePath)
	}
}
