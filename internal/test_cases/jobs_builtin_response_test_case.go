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

func (t JobsBuiltinResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) (err error) {
	defer func() {
		if err == nil && t.SuccessMessage != "" {
			logger.Successf("%s", t.SuccessMessage)
		}
	}()

	command := "jobs"

	if err := shell.SendCommand(command); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	// If we don't expect any items directly assert next prompt
	if len(t.ExpectedOutputItems) == 0 {
		return asserter.AssertWithPrompt()
	}

	for i, outputEntry := range t.ExpectedOutputItems {
		marker := convertJobMarkerToString(outputEntry.Marker)
		regexString := fmt.Sprintf(
			`^\[%d\]\s*%s\s+(?i)%s\s+(?-i)%s`,
			outputEntry.JobNumber,
			marker,
			regexp.QuoteMeta(outputEntry.Status),
			regexp.QuoteMeta(outputEntry.LaunchCommand),
		)

		// For 'running' jobs, bash displays the trailing & sign
		// This is optional since ZSH doesn't use this
		if strings.ToLower(outputEntry.Status) == "running" {
			regexString += "( &)?$"
		} else {
			regexString += "$"
		}

		regex := regexp.MustCompile(regexString)

		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput: fmt.Sprintf(
				"[%d]%s  %s                 %s",
				outputEntry.JobNumber, marker, outputEntry.Status, outputEntry.LaunchCommand,
			),
			FallbackPatterns: []*regexp.Regexp{regex},
		})

		shouldAssertWithPrompt := false

		if i == len(t.ExpectedOutputItems)-1 {
			shouldAssertWithPrompt = true
		}

		var err error

		if shouldAssertWithPrompt {
			err = asserter.AssertWithPrompt()
		} else {
			err = asserter.AssertWithoutPrompt()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func convertJobMarkerToString(jobMarker int) string {
	switch jobMarker {
	case UnmarkedJob:
		return " "
	case CurrentJob:
		return "+"
	case PreviousJob:
		return "-"
	}
	panic(fmt.Sprintf(
		"Codecrafters Internal Error: convertJobMarkerToString: Invalid job marker: %d",
		jobMarker,
	))
}
