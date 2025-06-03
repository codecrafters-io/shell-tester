// Stage HP4: History Persistence Test
// This test checks if the shell persists command history to a file when exiting.
// It runs some commands, exits the shell, and verifies the history file contents.

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
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

	// Step 4: Check history before exiting
	historyBefore := test_cases.HistoryTestCase{
		CommandsBeforeHistory: commandTestCases,
		SuccessMessage:        "✓ History before exiting is correct",
	}
	if err := historyBefore.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Step 5: Exit the shell
	shell.SendCommand("exit 0")
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ exit 0",
	})
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	// Read the output after exit command
	output := ""
	if screen := shell.GetScreenState(); len(screen) > 0 {
		output = strings.TrimSpace(screen[len(screen)-1][0])
	}
	// Allow both no output and 'exit' (like bash)
	if len(output) > 0 && output != "exit" {
		return fmt.Errorf("Expected no output or 'exit' after exit command, got %q", output)
	}

	asserter.LogRemainingOutput()

	// Step 6: Read the history file content
	historyContent, err := os.ReadFile(historyFile)
	if err != nil {
		return fmt.Errorf("failed to read history file: %v", err)
	}

	// Check if all commands are present in the history file
	historyStr := strings.TrimSpace(string(historyContent))
	historyLines := strings.Split(historyStr, "\n")

	utils.LogReadableFileContents(logger, historyStr, "History file content:", filepath.Base(historyFile))

	if len(historyLines) != len(commandTestCases)+2 {
		return fmt.Errorf("history file has %d lines, expected %d", len(historyLines), len(commandTestCases)+2)
	}
	logger.Successf("✓ History file has correct number of lines")

	// Verify new commands were written
	commandTestCases = append(commandTestCases, test_cases.CommandResponseTestCase{
		Command:        "history",
		ExpectedOutput: "history",
		SuccessMessage: "✓ Found command history in history file",
	})
	commandTestCases = append(commandTestCases, test_cases.CommandResponseTestCase{
		Command:        "exit 0",
		ExpectedOutput: "exit 0",
		SuccessMessage: "✓ Found command exit 0 in history file",
	})
	if len(historyLines) != len(commandTestCases) {
		return fmt.Errorf("history file has %d lines, expected %d", len(historyLines), len(commandTestCases))
	}
	for i, cmd := range commandTestCases {
		if historyLines[i] != cmd.Command {
			return fmt.Errorf("expected command %q at line %d, got %q", cmd.Command, i+1, historyLines[i])
		}
		logger.Successf("✓ Found command %q in history file", cmd.Command)
	}

	return logAndQuit(asserter, nil)
}
