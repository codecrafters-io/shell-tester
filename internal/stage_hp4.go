// Stage HP5: History Load on Startup Test
// This test checks if the shell loads command history from a file when starting up.
// It creates a history file with some commands, starts the shell, and verifies the commands are loaded.

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

func testHP4(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file with some initial commands
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+".txt")
	defer os.Remove(historyFile)

	// Create initial history file content
	nInitialCommands := random.RandomInt(2, 5)
	initialCommands := make([]string, nInitialCommands)
	for i := 0; i < nInitialCommands; i++ {
		initialCommands[i] = "echo " + strings.Join(random.RandomWords(random.RandomInt(2, 4)), " ")
	}
	initialContent := strings.Join(initialCommands, "\n") + "\n"
	if err := os.WriteFile(historyFile, []byte(initialContent), 0644); err != nil {
		return fmt.Errorf("failed to create history file: %v", err)
	}

	// Print the history file content
	historyContent, err := os.ReadFile(historyFile)
	if err != nil {
		panic("Codecrafters Internal Error: Failed to read history file")
	}
	utils.LogReadableFileContents(logger, string(historyContent), fmt.Sprintf("Writing commands to %s", historyFile), historyFile)

	// Set HISTFILE environment variable before starting the shell
	shell.Setenv("HISTFILE", historyFile)

	logger.Infof("export HISTFILE=%s", historyFile)
	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Step 2: Check history to verify initial commands are loaded
	historyTest := test_cases.HistoryTestCase{
		PreviousCommands: initialCommands,
		SuccessMessage:   "✓ History loaded from file is correct",
	}
	if err := historyTest.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
