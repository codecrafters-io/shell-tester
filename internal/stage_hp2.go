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
	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHP2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Create temporary history file paths
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+"_shell_history_test")
	defer os.Remove(historyFile) // Clean up the history file when done

	// Set HISTFILE to /dev/null before starting the shell
	shell.Setenv("HISTFILE", "/dev/null")

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Generate some random commands
	nCommands := 3
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
		SuccessMessage: "✓ History -w command executed",
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

	utils.LogReadableFileContents(logger, string(historyContent), "History file content:")

	// Check if all commands are present in the history file in order, allowing for number prefixes
	historyStr := string(historyContent)
	expectedCommands := make([]string, nCommands+1)
	for i, cmd := range commandTestCases {
		expectedCommands[i] = cmd.Command
	}
	historyLines := strings.Split(historyStr, "\n")
	expectedCommands[nCommands] = historyWriteCommand
	for i, cmd := range expectedCommands {
		if historyLines[i] != cmd {
			return fmt.Errorf("command %q not found in history file", cmd)
		}
		logger.Successf("✓ Found command %q in history file", cmd)
	}

	return logAndQuit(asserter, nil)
}
