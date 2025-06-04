// Stage HP3: History Append Test
// This test checks if the shell supports appending command history to a file using the history -a command.
// It runs some commands, appends them to history, and verifies the history file contents.

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

func testHP3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file with some initial content
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+".txt")

	// Create initial history file content (randomized, like in HP1/HP2)
	nInitialCommands := random.RandomInt(2, 4)
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

	utils.LogReadableFileContents(logger, initialContent, "Original history file content:", historyFile)

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

	// Step 3: Run commands
	for _, cmd := range commandTestCases {
		if err := cmd.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	// Step 4: Append history to file
	historyAppendCmd := "history -a " + historyFile
	historyAppendTest := test_cases.CommandWithNoRepsonseTestCase{
		Command:        historyAppendCmd,
		SuccessMessage: "✓ Ran history -a command",
	}
	if err := historyAppendTest.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Check if all commands are present in the history file
	commands := []string{}
	commands = append(commands, initialCommands...)
	for _, cmd := range commandTestCases {
		commands = append(commands, cmd.Command)
	}
	commands = append(commands, historyAppendCmd)

	if err := test_cases.AssertFileHasCommandsInOrder(logger, historyFile, commands); err != nil {
		return err
	}

	logger.Successf("✓ Found %d commands in history file", len(commands))

	// Run one echo command
	echoString := strings.Join(random.RandomWords(random.RandomInt(2, 4)), " ")
	echoTest := test_cases.CommandResponseTestCase{
		Command:        "echo " + echoString,
		ExpectedOutput: echoString,
		SuccessMessage: fmt.Sprintf("✓ Ran %s", "echo "+echoString),
	}

	if err := echoTest.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Run history -a again
	historyAppendTest = test_cases.CommandWithNoRepsonseTestCase{
		Command:        historyAppendCmd,
		SuccessMessage: "✓ Ran history -a command again",
	}
	if err := historyAppendTest.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Check if all commands are present in the history file
	commands = append(commands, echoTest.Command)
	commands = append(commands, historyAppendCmd)
	if err := test_cases.AssertFileHasCommandsInOrder(logger, historyFile, commands); err != nil {
		return err
	}

	logger.Successf("✓ Found %d commands in history file after running history -a command again", len(commands))

	return logAndQuit(asserter, nil)
}
