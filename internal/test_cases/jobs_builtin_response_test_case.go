package test_cases

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

const (
	UnmarkedJob = iota
	CurrentJob
	PreviousJob
)

type JobsBuiltinOutputEntry struct {
	// The job number value in the square brackets
	JobNumber int
	// Status: "Running", "Done", "Terminated", etc
	Status string
	// LaunchCommand: Command that was run and sent to the background without trailing &
	LaunchCommand string
	// Unmarked | Current | Previous
	Marker int
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
		switch outputEntry.Marker {
		case CurrentJob:
			marker = `\+`
		case PreviousJob:
			marker = `\-`
		}

		// This regex expects the following:
		// 1. Bracketed job number: Square bracket open, followed by an integer, followed by square bracket close
		// 2. An optional space (ZSH uses this optional space after bracketed job number)
		// 3. Job Marker (+/-/space)
		// 4. Whitespaces following the job marker
		// 5. Job status: "Done", "Running", etc (This is case insensitive: Complies with both bash and zsh)
		// 6. Followed by whitespace
		// 7. Followed by the launch command (case sensitive)
		regexString := fmt.Sprintf(
			`\[%d+\](\s)?%s\s+(?i)%s\s+(?-i)%s`,
			outputEntry.JobNumber,
			marker,
			regexp.QuoteMeta(outputEntry.Status),
			regexp.QuoteMeta(outputEntry.LaunchCommand),
		)

		// For 'running' jobs, bash displays the trailing & sign
		// This is optional since ZSH doesn't use this
		if strings.ToLower(outputEntry.Status) == "running" {
			regexString += "( &)?"
		}

		regex := regexp.MustCompile(regexString)

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
