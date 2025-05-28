// Stage HP2: History Write Test
// This test checks if the shell supports writing command history to a file using the history -w command.
// It sends a few commands, writes them to history, and verifies the history file contents.

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHP2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Create temporary history file paths
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+"_shell_history_test")
	overwriteFile := filepath.Join(os.TempDir(), random.RandomWord()+"_shell_history_test")
	defer os.Remove(historyFile)   // Clean up the history file when done
	defer os.Remove(overwriteFile) // Clean up the overwrite file when done

	// Generate some random commands
	nCommands := random.RandomInt(2, 5)
	commandTestCases := make([]test_cases.CommandResponseTestCase, nCommands)
	for i := 0; i < nCommands; i++ {
		randomWords := strings.Join(random.RandomWords(2), " ")
		commandTestCases[i] = test_cases.CommandResponseTestCase{
			Command:        "echo " + randomWords,
			ExpectedOutput: randomWords,
			SuccessMessage: "✓ Received expected response",
		}
	}

	// Execute initial commands
	for _, command := range commandTestCases {
		if err := command.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	// Write history to first file
	historyWriteCommand := "history -w " + historyFile
	logger.Infof("Writing history to file using command: %s", historyWriteCommand)

	historyTestCase := test_cases.CommandReflectionTestCase{
		Command:        historyWriteCommand,
		SuccessMessage: "✓ History written to file",
	}
	if err := historyTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Verify first history file contents
	logger.Infof("Verifying history file contents")
	historyContent, err := os.ReadFile(historyFile)
	if err != nil {
		logger.Errorf("Failed to read history file: %v", err)
		return err
	}

	// Check if all commands are present in the history file in order, allowing for number prefixes
	historyStr := string(historyContent)
	expectedCommands := make([]string, nCommands+1)
	for i, cmd := range commandTestCases {
		expectedCommands[i] = "echo " + cmd.ExpectedOutput
	}
	expectedCommands[nCommands] = historyWriteCommand
	for _, cmd := range expectedCommands {
		if !strings.Contains(historyStr, cmd) {
			logger.Errorf("Command %q not found in history file", cmd)
			return fmt.Errorf("command %q not found in history file", cmd)
		}
		logger.Successf("✓ Found command %q in history file", cmd)
	}

	// --- Test file overwrite behavior ---

	// Pre-populate the overwrite file with dummy content
	dummyContent := "this should be overwritten"
	if err := os.WriteFile(overwriteFile, []byte(dummyContent), 0644); err != nil {
		logger.Errorf("Failed to pre-populate overwrite file: %v", err)
		return err
	}

	// Run history -w to write to the overwrite file
	overwriteCommand := "history -w " + overwriteFile
	logger.Infof("Testing overwrite: writing history to a new file")
	historyTestCase2 := test_cases.CommandReflectionTestCase{
		Command:        overwriteCommand,
		SuccessMessage: "✓ History written to overwrite file",
	}
	if err := historyTestCase2.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Verify the overwrite file was written correctly
	overwriteContent, err := os.ReadFile(overwriteFile)
	if err != nil {
		logger.Errorf("Failed to read overwrite file: %v", err)
		return err
	}
	overwriteStr := string(overwriteContent)
	if strings.Contains(overwriteStr, dummyContent) {
		logger.Errorf("Overwrite file was not written correctly, still contains dummy content")
		return fmt.Errorf("overwrite file not written correctly")
	}
	for _, cmd := range expectedCommands {
		if !strings.Contains(overwriteStr, cmd) {
			logger.Errorf("Command %q not found in overwrite file", cmd)
			return fmt.Errorf("command %q not found in overwrite file", cmd)
		}
		logger.Successf("✓ Found command %q in overwrite file", cmd)
	}

	return logAndQuit(asserter, nil)
}
