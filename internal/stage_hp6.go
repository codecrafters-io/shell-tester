// Stage HP6: History Append on Exit Test
// This test checks if the shell appends command history to a file when exiting.
// It creates a history file with some initial content, runs commands, exits the shell,
// and verifies the history file contents are appended.

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

func testHP6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file with some initial content
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+".txt")
	defer os.Remove(historyFile)

	// Create initial history file content
	nInitialCommands := random.RandomInt(2, 4)
	initialCommands := make([]string, nInitialCommands)
	for i := 0; i < nInitialCommands; i++ {
		initialCommands[i] = "echo " + strings.Join(random.RandomWords(random.RandomInt(2, 4)), " ")
	}
	initialContent := strings.Join(initialCommands, "\n") + "\n"
	if err := os.WriteFile(historyFile, []byte(initialContent), 0644); err != nil {
		return fmt.Errorf("failed to create history file: %v", err)
	}

	utils.LogReadableFileContents(logger, initialContent, "Original history file content:", historyFile)

	// Set HISTFILE environment variable before starting the shell
	shell.Setenv("HISTFILE", historyFile)

	logger.Infof("export HISTFILE=%s", historyFile)

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

	previousCommands := []string{}
	previousCommands = append(previousCommands, initialCommands...)
	for _, command := range commandTestCases {
		previousCommands = append(previousCommands, command.Command)
	}

	for _, command := range commandTestCases {
		if err := command.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	// Step 3: Check history before exiting
	historyBefore := test_cases.HistoryTestCase{
		PreviousCommands: previousCommands,
		SuccessMessage:   "✓ History before exiting is correct",
	}
	if err := historyBefore.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Step 4: Exit the shell
	exitTest := test_cases.ExitTestCase{
		Command:          "exit 0",
		ExpectedExitCode: 0,
	}
	if err := exitTest.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Step 5: Check history file contents
	commands := []string{}
	commands = append(commands, previousCommands...)
	commands = append(commands, "history")
	commands = append(commands, "exit 0")
	if err := test_cases.AssertHistoryFileHasCommands(logger, historyFile, commands); err != nil {
		return err
	}

	logger.Successf("✓ Found %d commands in history file", len(commands))

	return logAndQuit(asserter, nil)
}
