package test_cases

import (
	"regexp"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type CommandResponseWithReapedJobsTestCase struct {
	// Command is the command to send to the shell
	Command string

	// ExpectedCommandOutput is the expected output string to match against
	ExpectedCommandOutput string

	// FallbackPatterns is a list of regex patterns to match against
	FallbackPatterns []*regexp.Regexp

	// SuccessMessage is the message to log in case of success
	SuccessMessage string

	// ExpectedReapedJobEntries is the list of entries expected after the command output appears
	ExpectedReapedJobEntries []*BackgroundJobStatusEntry
}

func (t CommandResponseWithReapedJobsTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	allSingleLinesAssertion := []assertions.SingleLineAssertion{}

	// Add assertion for command's actual output first
	allSingleLinesAssertion = append(allSingleLinesAssertion, assertions.SingleLineAssertion{
		ExpectedOutput:   t.ExpectedCommandOutput,
		FallbackPatterns: t.FallbackPatterns,
	})

	// Add assertion for reaped job entries now
	for _, expectedOutputEntry := range t.ExpectedReapedJobEntries {
		expectedOutput, regexPattern := expectedOutputEntry.ExpectedOutputAndRegex()

		allSingleLinesAssertion = append(allSingleLinesAssertion, assertions.SingleLineAssertion{
			ExpectedOutput:   expectedOutput,
			FallbackPatterns: []*regexp.Regexp{regexPattern},
		})
	}

	// A small delay to ensure that the grep process has exitted
	time.Sleep(time.Millisecond)

	commandWithMultilineResponseTestCase := CommandWithMultilineResponseTestCase{
		Command: t.Command,
		MultiLineAssertion: assertions.MultiLineAssertion{
			SingleLineAssertions: allSingleLinesAssertion,
		},
		SuccessMessage: t.SuccessMessage,
	}

	if err := commandWithMultilineResponseTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}
