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
	ExpectedOutputEntries []JobsBuiltinOutputEntry
	SuccessMessage        string
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

	// In case of no output entries, assert only the command reflection
	if len(t.ExpectedOutputEntries) == 0 {
		return asserter.AssertWithoutPrompt()
	}

	for i, outputEntry := range t.ExpectedOutputEntries {
		marker := convertJobMarkerToString(outputEntry.Marker)
		// This regex aims to match lines like: [1]+  Running                 sleep 5 &
		regexString := fmt.Sprintf(
			`^\[%d\]\s*%s\s+(?i)%s\s+(?-i)%s`,
			outputEntry.JobNumber,
			regexp.QuoteMeta(marker),
			regexp.QuoteMeta(outputEntry.Status),
			regexp.QuoteMeta(outputEntry.LaunchCommand),
		)

		// For 'running' jobs, bash displays the trailing & sign
		// This is optional since ZSH doesn't use this
		if outputEntry.Status == "Running" {
			regexString += "( &)?$"
		} else {
			regexString += "$"
		}

		regex := regexp.MustCompile(regexString)

		expectedOutput := fmt.Sprintf(
			"[%d]%s  %s                 %s",
			outputEntry.JobNumber, marker, outputEntry.Status, outputEntry.LaunchCommand,
		)

		if outputEntry.Status == "Running" {
			expectedOutput += " &"
		}

		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput:   expectedOutput,
			FallbackPatterns: []*regexp.Regexp{regex},
		})

		assertWithPrompt := false

		if i == len(t.ExpectedOutputEntries)-1 {
			assertWithPrompt = true
		}

		var err error

		if assertWithPrompt {
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
