// Stage HP3: History Append Test
// This test checks if the shell supports appending command history to a file using the history -a command.
// It runs some commands, appends them to history, and verifies the history file contents.

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	utils.LogReadableFileContents(logger, initialContent, "Original history file content:", filepath.Base(historyFile))

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
	historyBeforeTest := test_cases.HistoryTestCase{
		CommandsBeforeHistory: commandTestCases,
		SuccessMessage:        "✓ History before appending is correct",
	}
	if err := historyBeforeTest.Run(asserter, shell, logger); err != nil {
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

	// Read the updated history file content
	historyContent, err := os.ReadFile(historyFile)
	if err != nil {
		logger.Errorf("Failed to read history file: %v", err)
		return err
	}

	utils.LogReadableFileContents(logger, string(historyContent), "History file content after appending:", filepath.Base(historyFile))

	// Check if all commands are present in the history file
	historyStr := strings.TrimSpace(string(historyContent))
	historyLines := strings.Split(historyStr, "\n")
	if len(historyLines) != len(initialCommands)+len(commandTestCases)+2 {
		return fmt.Errorf("history file has %d lines, expected %d", len(historyLines), len(initialCommands)+len(commandTestCases)+2)
	}
	logger.Successf("✓ History file has correct number of lines")
	i := 0
	// First verify initial commands are still present
	for _, cmd := range initialCommands {
		if historyLines[i] != cmd {
			return fmt.Errorf("expected initial command %q at line %d, got %q", cmd, i+1, historyLines[i])
		}
		logger.Successf("✓ Found initial command %q in history file", cmd)
		i++
	}

	// Then verify new commands were appended
	for _, cmd := range commandTestCases {
		if historyLines[i] != cmd.Command {
			return fmt.Errorf("expected new command %q at line %d, got %q", cmd.Command, i+1, historyLines[i])
		}
		logger.Successf("✓ Found new command %q in history file", cmd.Command)
		i++
	}

	// Verify history command itself is present
	if historyLines[i] != "history" {
		return fmt.Errorf("history command not found in history file")
	}
	logger.Successf("✓ Found history command in history file")
	i++

	// Verify history -a command itself is present
	if historyLines[i] != historyAppendCmd {
		return fmt.Errorf("expected history append command %q at line %d, got %q", historyAppendCmd, i+1, historyLines[i])
	}
	logger.Successf("✓ Found history append command in history file")

	// Get initial counts of commands
	initialCounts := make(map[string]int)
	for _, cmd := range initialCommands {
		initialCounts[cmd] = strings.Count(historyStr, cmd)
	}
	for _, cmd := range commandTestCases {
		initialCounts[cmd.Command] = strings.Count(historyStr, cmd.Command)
	}

	// Run history -a again
	historyAppendTest = test_cases.CommandReflectionTestCase{
		Command:        historyAppendCmd,
		SuccessMessage: "✓ history -a command executed",
	}
	if err := historyAppendTest.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// Read the updated history file content again
	historyContent, err = os.ReadFile(historyFile)
	if err != nil {
		logger.Errorf("Failed to read history file: %v", err)
		return err
	}
	historyStr = string(historyContent)

	utils.LogReadableFileContents(logger, string(historyContent), "History file content after second append:", filepath.Base(historyFile))

	// Verify counts haven't increased (excluding history -a command)

	// Sort the commands to ensure consistent order (for testing)
	var cmds []string
	for cmd := range initialCounts {
		if cmd != historyAppendCmd {
			cmds = append(cmds, cmd)
		}
	}
	sort.Strings(cmds)
	for _, cmd := range cmds {
		initialCount := initialCounts[cmd]
		newCount := strings.Count(historyStr, cmd)
		if newCount > initialCount {
			logger.Errorf("Command %q appears %d times in history file (was %d)", cmd, newCount, initialCount)
			return fmt.Errorf("command %q appears %d times in history file (was %d)", cmd, newCount, initialCount)
		}
		logger.Successf("✓ Command %q count preserved (%d)", cmd, initialCount)
	}

	return logAndQuit(asserter, nil)
}
