// Stage HP4: History Persistence Test
// This test checks if the shell persists command history to a file when exiting.
// It runs some commands, exits the shell, and verifies the history file contents.

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHP4(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+"_shell_history_test")
	defer os.Remove(historyFile)

	// Create the file
	if err := os.WriteFile(historyFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create history file: %v", err)
	}

	// Set HISTFILE environment variable before starting the shell
	shell.Setenv("HISTFILE", historyFile)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Step 3: Run some commands in the shell
	nShellCommands := 3 // Fixed number of shell commands
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

	// Step 4: Check history before exiting
	historyBefore := test_cases.HistoryTestCase{
		SuccessMessage:        "✓ History before exiting is correct",
		CommandsBeforeHistory: commandTestCases,
	}
	if err := historyBefore.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Step 5: Exit the shell
	shell.SendCommand("exit 0")
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ exit 0",
	})
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "exit",
		StayOnSameLine: true,
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^$`),
			regexp.MustCompile(`^exit$`),
		},
	})
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Add a small delay to ensure file is written
	time.Sleep(100 * time.Millisecond)

	// Step 6: Read the history file content
	historyContent, err := os.ReadFile(historyFile)
	if err != nil {
		logger.Errorf("Failed to read history file: %v", err)
		return err
	}

	// Check if all commands are present in the history file
	historyStr := string(historyContent)

	// print the history file content
	logger.Debugf("History file content: \n%s", historyStr)

	// Verify new commands were written
	for _, cmd := range commandTestCases {
		if !strings.Contains(historyStr, cmd.Command) {
			logger.Errorf("Command %q not found in history file", cmd.Command)
			return fmt.Errorf("command %q not found in history file", cmd.Command)
		}
		logger.Successf("✓ Found command %q in history file", cmd.Command)
	}

	return logAndQuit(asserter, nil)
}
