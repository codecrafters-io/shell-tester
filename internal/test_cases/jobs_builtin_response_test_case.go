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

// jobsBuiltinOutputLineRegex matches a single jobs output line and captures:
// 1. job id
// 2. optional whitespace between job id and marker (+/-/space)
// 3. marker (+/-/space)
// 4. whitespaces
// 5. Status (Could be a single non whitespace like "Done", "Running", or could be sth like "1 Exit"; This last case we'll be used in future extension
// 6. Whitespaces
// 7. Launch command
// 8. Optional & sign that bash uses
var jobsBuiltinOutputLineRegex = regexp.MustCompile(`^\[(\d+)\]\s*([\+\-\s])\s+(\S+( )?\S+)\s+(.*)( &)?$`)

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

	// In case of no output entries, assert only the command reflection
	if len(t.ExpectedOutputItems) == 0 {
		return asserter.AssertWithoutPrompt()
	}

	for i, outputEntry := range t.ExpectedOutputItems {
		asserter.AddAssertion(assertions.SingleLineRegexAssertion{
			ExpectedRegexPatterns: []*regexp.Regexp{jobsBuiltinOutputLineRegex},
		})

		assertWithPrompt := false

		if i == len(t.ExpectedOutputItems)-1 {
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

		outputLine := asserter.Shell.GetScreenState().GetRow(asserter.GetLastLoggedRowIndex())
		outputText := outputLine.String()

		if err := validateJobsOutputLineWithCaptures(outputText, outputEntry); err != nil {
			return err
		}
	}

	return nil
}

// validateJobsOutputLineWithCaptures parses the jobs output line with capture groups and compares
// each value to the expected entry. Returns a descriptive error on first mismatch.
// Panics with a Codecrafters Internal Error if the line does not match the expected format.
func validateJobsOutputLineWithCaptures(outputText string, expectedEntry JobsBuiltinOutputEntry) error {
	submatches := jobsBuiltinOutputLineRegex.FindStringSubmatch(outputText)

	if len(submatches) < 6 {
		panic("Codecrafters Internal Error: Shouldn't be here - Could not parse jobs output line")
	}

	capturedJobIDStr := submatches[1]
	capturedMarkerStr := submatches[2]
	capturedStatusStr := submatches[3]
	capturedLaunchCommandStr := submatches[5]

	capturedMarker := parsedMarkerToMarkerConstant(capturedMarkerStr)

	normalisedCapturedLaunchCommand := strings.TrimSuffix(strings.TrimSpace(capturedLaunchCommandStr), "&")
	normalisedCapturedLaunchCommand = strings.TrimSpace(normalisedCapturedLaunchCommand)
	expectedLaunchCommandNormalised := strings.TrimSpace(expectedEntry.LaunchCommand)

	if capturedJobIDStr != fmt.Sprintf("%d", expectedEntry.JobNumber) {
		return fmt.Errorf("Job number mismatch: expected %d, got %s", expectedEntry.JobNumber, capturedJobIDStr)
	}

	if capturedMarker != expectedEntry.Marker {
		return fmt.Errorf("Marker mismatch: expected %s, got %s",
			markerConstantToDisplay(expectedEntry.Marker), markerConstantToDisplay(capturedMarker))
	}

	if !strings.EqualFold(capturedStatusStr, expectedEntry.Status) {
		return fmt.Errorf("Status mismatch: expected %q, got %q", expectedEntry.Status, capturedStatusStr)
	}

	if normalisedCapturedLaunchCommand != expectedLaunchCommandNormalised {
		return fmt.Errorf("Launch command mismatch: expected %q, got %q",
			expectedEntry.LaunchCommand, capturedLaunchCommandStr)
	}

	return nil
}

func parsedMarkerToMarkerConstant(markerStr string) int {
	switch strings.TrimSpace(markerStr) {
	case "+":
		return CurrentJob
	case "-":
		return PreviousJob
	default:
		return UnmarkedJob
	}
}

func markerConstantToDisplay(marker int) string {
	switch marker {
	case CurrentJob:
		return "\"+\" (current)"
	case PreviousJob:
		return "\"-\" (previous)"
	default:
		return "\" \" (unmarked)"
	}
}
