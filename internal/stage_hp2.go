// Stage HP2: History Write Test
// This test checks if the shell supports writing command history to a file using the history -w command.
// It sends a few commands, writes them to history, and verifies the history file contents.

package internal

import (
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

	// Create temporary history file paths
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+".txt")
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

	// Write history to file
	historyWriteCommand := "history -w " + historyFile
	logger.Infof("Writing history to file using command: %s", historyWriteCommand)

	historyTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        historyWriteCommand,
		SuccessMessage: "✓ Ran history -w command",
	}
	if err := historyTestCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	// Verify first history file contents
	commands := []string{}
	for _, command := range commandTestCases {
		commands = append(commands, command.Command)
	}
	commands = append(commands, historyWriteCommand)

	if err := test_cases.AssertHistoryFileHasCommands(logger, historyFile, commands); err != nil {
		return err
	}

	logger.Successf("✓ Found %d commands in history file", len(commands))

	return logAndQuit(asserter, nil)
}
