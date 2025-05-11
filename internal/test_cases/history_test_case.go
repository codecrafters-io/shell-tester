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
	// Run all commands before history
	for _, cmdPair := range t.CommandsBeforeHistory {
		if err := shell.SendCommand(cmdPair.Command); err != nil {
			return fmt.Errorf("Error sending command to shell: %v", err)
		}

		// Check command reflection
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: fmt.Sprintf("$ %s", cmdPair.Command),
		})

		// Check command output if expected
		if cmdPair.ExpectedOutput != "" {
			asserter.AddAssertion(assertions.SingleLineAssertion{
				ExpectedOutput: cmdPair.ExpectedOutput,
			})
		}

		if err := asserter.AssertWithPrompt(); err != nil {
			return err
		}
	}

	// Now run the history command
	if err := shell.SendCommand("history"); err != nil {
		return fmt.Errorf("Error sending history command to shell: %v", err)
	}

	// Check if the history command is present in the output
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "$ history",
	})

	// Add assertions for each command in history
	for _, cmdPair := range t.CommandsBeforeHistory {
		// Add assertion for each command in history
		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: "",
			FallbackPatterns: []*regexp.Regexp{
				regexp.MustCompile(fmt.Sprintf(`^\s*\d+\s+%s$`, regexp.QuoteMeta(cmdPair.Command))),
			},
		})
	}

	// Add assertion for the history command itself
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: "",
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*\d+\s+history$`),
		},
	})

	// Log the history output before asserting
	logger.Infof("History output:")
	asserter.LogRemainingOutput()

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
