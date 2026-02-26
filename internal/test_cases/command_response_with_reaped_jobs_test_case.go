package test_cases

import (
	"fmt"
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
		expectedJobMarkerString := convertJobMarkerToString(expectedOutputEntry.Marker)

		expectedOutput := fmt.Sprintf(
			"[%d]%s  %s                 %s",
			expectedOutputEntry.JobNumber, expectedJobMarkerString, expectedOutputEntry.Status, expectedOutputEntry.LaunchCommand,
		)

		// This regex aims to match lines like: [1]+  Done                 sleep 5 &
		regexString := fmt.Sprintf(
			`^\[%d\]\s*%s\s+(?i)%s\s+(?-i)%s$`,
			expectedOutputEntry.JobNumber,
			regexp.QuoteMeta(expectedJobMarkerString),
			regexp.QuoteMeta(expectedOutputEntry.Status),
			regexp.QuoteMeta(expectedOutputEntry.LaunchCommand),
		)

		allSingleLinesAssertion = append(allSingleLinesAssertion, assertions.SingleLineAssertion{
			ExpectedOutput:   expectedOutput,
			FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(regexString)},
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
