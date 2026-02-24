package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

const (
	UnlabeledJob = iota
	CurrentJob
	PreviousJob
)

type JobsBuiltinOutputEntry struct {
	// The job number value in the square brackets
	JobNumber int
	// Status: "Running", "Done", etc
	Status string
	// LaunchCommand: Command that was run and sent to the background without trailing &
	LaunchCommand string
	// Unlabeled | Current | Previous
	Label int
}

type JobsBuiltinResponseTestCase struct {
	ExpectedOutputItems []JobsBuiltinOutputEntry
	SuccessMessage      string
}

func (t JobsBuiltinResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	command := "jobs"

	if err := shell.SendCommand(command); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	allLinesAssertions := []assertions.SingleLineAssertion{}

	for _, outputEntry := range t.ExpectedOutputItems {
		marker := `\s`
		switch outputEntry.Label {
		case CurrentJob:
			marker = `\+`
		case PreviousJob:
			marker = `\-`
		}

		// TODO: Remove after PR review: The regex here complies with the Bash's implementation of 'jobs'
		// Should I add the pattern compatible with ZSH's output as well?
		regex := regexp.MustCompile(fmt.Sprintf(
			`\[%d\]%s\s+%s\s+%s &`,
			outputEntry.JobNumber,
			marker,
			regexp.QuoteMeta(outputEntry.Status),
			regexp.QuoteMeta(outputEntry.LaunchCommand)),
		)

		allLinesAssertions = append(allLinesAssertions, assertions.SingleLineAssertion{
			FallbackPatterns: []*regexp.Regexp{regex},
		})
	}

	asserter.AddAssertion(&assertions.MultiLineAssertion{
		SingleLineAssertions: allLinesAssertions,
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
