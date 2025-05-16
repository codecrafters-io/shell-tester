package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// CommandOutputPair represents a command and its expected output
type CommandOutputPair struct {
	Command        string
	ExpectedOutput string
}

// HistoryTestCase is a test case that:
// Sends the history command to the shell
// Verifies that the command is printed to the screen
// Verifies that the history output shows the command in the expected format
type HistoryTestCase struct {
	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// CommandsBeforeHistory is a list of commands to execute before running history
	CommandsBeforeHistory []CommandOutputPair
}

func (t HistoryTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	for _, cmdPair := range t.CommandsBeforeHistory {
		if err := shell.SendCommand(cmdPair.Command); err != nil {
			return fmt.Errorf("failed to send command: %v", err)
		}

		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: fmt.Sprintf("$ %s", cmdPair.Command),
		})

		if cmdPair.ExpectedOutput != "" {
			asserter.AddAssertion(assertions.SingleLineAssertion{
				ExpectedOutput: cmdPair.ExpectedOutput,
			})
		}

		if err := asserter.AssertWithPrompt(); err != nil {
			return err
		}
	}

	if err := shell.SendCommand("history"); err != nil {
		return fmt.Errorf("failed to send history command: %v", err)
	}

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ history",
	})

	for i, cmdPair := range t.CommandsBeforeHistory {
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: fmt.Sprintf("    %d  %s", i+1, cmdPair.Command),
			FallbackPatterns: []*regexp.Regexp{
				regexp.MustCompile(fmt.Sprintf(`^    \d+\s+%s$`, regexp.QuoteMeta(cmdPair.Command))),
			},
		})
	}

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: fmt.Sprintf("    %d  history", len(t.CommandsBeforeHistory)+1),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^    \d+\s+history$`),
		},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
