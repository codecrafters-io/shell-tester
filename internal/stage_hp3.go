// Stage HP3: History Append Test
// This test checks if the shell supports appending command history to a file using the history -a command.
// It runs some commands, appends them to history, and verifies the history file contents.

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHP3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file with some initial content
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+"_shell_history_test")

	// Create initial history file content (randomized, like in HP1/HP2)
	nInitialCommands := random.RandomInt(2, 5)
	initialCommands := make([]string, nInitialCommands)
	for i := 0; i < nInitialCommands; i++ {
		initialCommands[i] = "echo " + strings.Join(random.RandomWords(random.RandomInt(2, 4)), " ")
	}
	initialContent := strings.Join(initialCommands, "\n") + "\n"
	if err := os.WriteFile(historyFile, []byte(initialContent), 0666); err != nil {
		logger.Errorf("Failed to create history file: %v", err)
		return err
	}
	defer os.Remove(historyFile)

	// Set HISTFILE to /dev/null before starting the shell
	shell.Setenv("HISTFILE", "/dev/null")

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Step 2: Run some commands in the shell
	nShellCommands := random.RandomInt(2, 4)
	commandTestCases := make([]test_cases.CommandResponseTestCase, nShellCommands)
	for i := 0; i < nShellCommands; i++ {
		cmdWords := random.RandomWords(random.RandomInt(2, 4))
		cmd := "echo " + strings.Join(cmdWords, " ")
		commandTestCases[i] = test_cases.CommandResponseTestCase{
			Command:        cmd,
			ExpectedOutput: strings.Join(cmdWords, " "),
			SuccessMessage: fmt.Sprintf("✓ Ran %s", cmd),
		}
	}

	// Step 3: Check history before appending
	historyBefore := test_cases.HistoryTestCase{
		SuccessMessage:        "✓ History before appending is correct",
		CommandsBeforeHistory: commandTestCases,
	}
	if err := historyBefore.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Step 4: Append history to file
	historyAppendCmd := "history -a " + historyFile
	historyAppendTest := test_cases.CommandReflectionTestCase{
		Command:        historyAppendCmd,
		SuccessMessage: "✓ History appended to file",
	}
	if err := historyAppendTest.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Add a small delay to ensure file is written
	time.Sleep(100 * time.Millisecond)

	// Read the updated history file content
	historyContent, err := os.ReadFile(historyFile)
	if err != nil {
		logger.Errorf("Failed to read history file: %v", err)
		return err
	}

	// Check if all commands are present in the history file
	historyStr := string(historyContent)

	// First verify initial commands are still present
	for _, cmd := range initialCommands {
		if !strings.Contains(historyStr, cmd) {
			logger.Errorf("Initial command %q not found in history file", cmd)
			return fmt.Errorf("initial command %q not found in history file", cmd)
		}
		logger.Successf("✓ Found initial command %q in history file", cmd)
	}

	// Then verify new commands were appended
	for _, cmd := range commandTestCases {
		if !strings.Contains(historyStr, cmd.Command) {
			logger.Errorf("New command %q not found in history file", cmd.Command)
			return fmt.Errorf("new command %q not found in history file", cmd.Command)
		}
		logger.Successf("✓ Found new command %q in history file", cmd.Command)
	}

	// Verify history -a command itself is present
	if !strings.Contains(historyStr, historyAppendCmd) {
		logger.Errorf("History append command %q not found in history file", historyAppendCmd)
		return fmt.Errorf("history append command %q not found in history file", historyAppendCmd)
	}
	logger.Successf("✓ Found history append command in history file")

	return logAndQuit(asserter, nil)
}
