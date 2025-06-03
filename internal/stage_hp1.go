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
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHP1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	// Step 1: Create a temporary history file with some commands
	historyFile := filepath.Join(os.TempDir(), random.RandomWord()+"_shell_history_test")
	defer os.Remove(historyFile)

	// Set HISTFILE to /dev/null before starting the shell
	shell.Setenv("HISTFILE", "/dev/null")

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Step 1: Create history file with random commands
	nFileCommands := 3 // Fixed number of initial commands
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
	logger.Debugf("History file content: \n%s", string(content))

	// Step 2: Run some commands in the shell
	nShellCommands := 2 // Fixed number of shell commands
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

	// Step 3: Check history before loading from file
	historyBefore := test_cases.HistoryTestCase{
		SuccessMessage:        "✓ History before loading file is correct",
		CommandsBeforeHistory: commandTestCases,
	}
	if err := historyBefore.Run(asserter, shell, logger); err != nil {
		return err
	}

	// Step 4: Load history from file
	historyLoadCmd := "history -r " + historyFile
	historyLoadTest := test_cases.CommandReflectionTestCase{
		Command:        historyLoadCmd,
		SuccessMessage: "✓ Loaded history from file",
	}
	if err := historyLoadTest.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Step 5: Check history after loading from file (should include previous + loaded)
	historyAfter := test_cases.HistoryPersistenceTestCase{
		PreviousCommands:          commandTestCases,
		FilePath:                  historyFile,
		SuccessMessage:            "✓ History after loading file is correct",
		WasHistoryCommandExecuted: true,
		ExpectHistoryRCommand:     true,
	}
	if err := historyAfter.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
