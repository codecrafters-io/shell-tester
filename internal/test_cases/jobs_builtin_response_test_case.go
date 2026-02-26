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

type BackgroundJobStatusEntry struct {
	// The job number value in the square brackets
	JobNumber int
	// Status: "Running", "Done", "Terminated", "1 Exit", etc
	Status string
	// LaunchCommand: Command that was run and sent to the background without trailing &
	LaunchCommand string
	// Unmarked | Current | Previous
	Marker int
}

// ExpectedOutputAndRegex returns the expected output string and regex pattern for this job entry.
func (e BackgroundJobStatusEntry) ExpectedOutputAndRegex() (string, *regexp.Regexp) {
	expectedJobMarkerString := convertJobMarkerToString(e.Marker)

	// This regex aims to match lines like: [1]+  Running                 sleep 5 &
	regexString := fmt.Sprintf(
		`^\[%d\]\s*%s\s+(?i)%s\s+(?-i)%s`,
		e.JobNumber,
		regexp.QuoteMeta(expectedJobMarkerString),
		regexp.QuoteMeta(e.Status),
		regexp.QuoteMeta(e.LaunchCommand),
	)

	// For 'Running' jobs, bash displays the trailing & sign
	// Users shall comply with bash for consistency (Ensured this by appending this to expected output)
	// But this should be optional since ZSH doesn't use this
	if e.Status == "Running" {
		regexString += "( &)?$"
	} else {
		regexString += "$"
	}

	expectedOutput := fmt.Sprintf(
		"[%d]%s  %s                 %s",
		e.JobNumber, expectedJobMarkerString, e.Status, e.LaunchCommand,
	)

	// For 'Running' jobs, the trailing sign is expected
	if e.Status == "Running" {
		expectedOutput += " &"
	}

	return expectedOutput, regexp.MustCompile(regexString)
}

type JobsBuiltinResponseTestCase struct {
	ExpectedOutputEntries []BackgroundJobStatusEntry
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
		return asserter.AssertWithPrompt()
	}

	for i, expectedOutputEntry := range t.ExpectedOutputEntries {
		expectedOutput, regexPattern := expectedOutputEntry.ExpectedOutputAndRegex()

		asserter.AddAssertion(assertions.SingleLineAssertion{
			ExpectedOutput:   expectedOutput,
			FallbackPatterns: []*regexp.Regexp{regexPattern},
		})

		assertWithPrompt := false
		var err error

		if i == len(t.ExpectedOutputEntries)-1 {
			assertWithPrompt = true
		}

		// Assert with prompt on last entry
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
