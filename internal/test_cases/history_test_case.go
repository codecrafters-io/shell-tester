package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// HistoryTestCase is a test case that:
// Sends the history command to the shell
// Verifies that the command is printed to the screen
// Verifies that the history output shows the command in the expected format
type HistoryTestCase struct {
	// HistoryOffset specifies the starting line number offset for the PreviousCommands
	// This accounts for commands that were executed before the PreviousCommands in the same session
	HistoryOffset int

	// LastNCommands specifies how many of the most recent commands to check in history
	// If not set, all commands will be checked
	LastNCommands int

	// PreviousCommands is a list of previous commands expected to be in the history list
	PreviousCommands []string

	// SuccessMessage is the message to log in case of success
	SuccessMessage string
}

func (t HistoryTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	historyCommand := "history"
	if t.LastNCommands > 0 {
		historyCommand = fmt.Sprintf("history %d", t.LastNCommands)
	}
	historyReflectionTest := CommandWithNoResponseTestCase{
		Command:             historyCommand,
		SuccessMessage:      "✓ Ran history command",
		SkipPromptAssertion: true,
	}
	if err := historyReflectionTest.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	// Calculate which commands to check based on LastNCommands
	startIdx := 0
	if t.LastNCommands > 0 && t.LastNCommands < len(t.PreviousCommands) {
		startIdx = len(t.PreviousCommands) - t.LastNCommands + 1
	}

	// Check only the specified number of most recent commands
	for i, command := range t.PreviousCommands[startIdx:] {
		expectedLineNumber := t.HistoryOffset + startIdx + i + 1
		asserter.AddAssertion(&assertions.SingleLineAssertion{
			ExpectedOutput: fmt.Sprintf("    %d  %s", expectedLineNumber, command),
			FallbackPatterns: []*regexp.Regexp{
				regexp.MustCompile(fmt.Sprintf(`^\s*%d\s+%s$`, expectedLineNumber, regexp.QuoteMeta(command))),
				regexp.MustCompile(fmt.Sprintf(`^\s*%d\s+%s$`, expectedLineNumber-1, regexp.QuoteMeta(command))), // 0-based for ash
			},
		})
	}

	// Add assertion for the history command itself
	expectedHistoryLineNumber := t.HistoryOffset + len(t.PreviousCommands) + 1
	asserter.AddAssertion(&assertions.SingleLineAssertion{
		ExpectedOutput: fmt.Sprintf("    %d  %s", expectedHistoryLineNumber, historyCommand),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^\s*%d\s+%s$`, expectedHistoryLineNumber, regexp.QuoteMeta(historyCommand))),
			regexp.MustCompile(fmt.Sprintf(`^\s*%d\s+%s$`, expectedHistoryLineNumber-1, regexp.QuoteMeta(historyCommand))), // 0-based for ash
		},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
