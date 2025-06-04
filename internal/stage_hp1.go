// Stage HP1: History Read Test
// This test checks if the shell supports reading command history from a file using the history -r command.
// It creates a history file with some commands, loads it into the shell, and verifies the history contents.

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

func testHP1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file with some commands
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+".txt")
	defer os.Remove(historyFile)

	// Set HISTFILE to /dev/null before starting the shell
	shell.Setenv("HISTFILE", "/dev/null")

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Step 1: Create history file with random commands
	nFileCommands := random.RandomInt(2, 4)
	fileCommands := make([]string, nFileCommands)
	for i := 0; i < nFileCommands; i++ {
		fileCommands[i] = "echo " + strings.Join(random.RandomWords(random.RandomInt(2, 4)), " ")
	}
	historyFileContent := strings.Join(fileCommands, "\n") + "\n"
	if err := os.WriteFile(historyFile, []byte(historyFileContent), 0644); err != nil {
		logger.Errorf("Failed to write history file: %v", err)
		return err
	}

	// print the history file content
	content, err := os.ReadFile(historyFile)
	if err != nil {
		logger.Errorf("Failed to read history file: %v", err)
		return err
	}
	utils.LogReadableFileContents(logger, string(content), fmt.Sprintf("Writing contents to %s", historyFile), historyFile)

	// Step 2: Load history from file
	historyLoadTest := test_cases.CommandWithNoResponseTestCase{
		Command:        fmt.Sprintf("history -r %s", historyFile),
		SuccessMessage: "✓ Ran history -r command",
	}
	if err := historyLoadTest.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Step 3: Check history after loading from file (should include previous + loaded)
	fileCommandsToPass := []string{}
	fileCommandsToPass = append(fileCommandsToPass, fmt.Sprintf("history -r %s", historyFile))
	fileCommandsToPass = append(fileCommandsToPass, fileCommands...)

	afterLoadHistoryTest := test_cases.HistoryTestCase{
		PreviousCommands: fileCommandsToPass,
		SuccessMessage:   "✓ History after loading file is correct",
	}
	if err := afterLoadHistoryTest.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
