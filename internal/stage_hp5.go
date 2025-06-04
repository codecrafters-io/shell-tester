// Stage HP4: History Persistence Test
// This test checks if the shell persists command history to a file when exiting.
// It runs some commands, exits the shell, and verifies the history file contents.

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

func testHP5(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+".txt")
	defer os.Remove(historyFile)

	// Create the file
	if err := os.WriteFile(historyFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create history file: %v", err)
	}

	// Set HISTFILE environment variable before starting the shell
	shell.Setenv("HISTFILE", historyFile)

	logger.Infof("export HISTFILE=%s", historyFile)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Step 3: Run some commands in the shell
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
	for _, command := range commandTestCases {
		previousCommands = append(previousCommands, command.Command)
	}

	for _, command := range commandTestCases {
		if err := command.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	// Step 4: Check history before exiting
	historyBefore := test_cases.HistoryTestCase{
		PreviousCommands: previousCommands,
		SuccessMessage:   "✓ History before exiting is correct",
	}
	if err := historyBefore.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Step 5: Exit the shell
	exitTest := test_cases.ExitTestCase{
		Command:          "exit 0",
		ExpectedExitCode: 0,
	}
	if err := exitTest.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Step 6: Check history file contents
	commands := make([]string, len(commandTestCases))
	for i, cmd := range commandTestCases {
		commands[i] = cmd.Command
	}
	commands = append(commands, "history")
	commands = append(commands, "exit 0")
	if err := test_cases.AssertFileHasCommandsInOrder(logger, historyFile, commands); err != nil {
		return err
	}

	logger.Successf("✓ Found %d commands in history file", len(commands))

	return logAndQuit(asserter, nil)
}
