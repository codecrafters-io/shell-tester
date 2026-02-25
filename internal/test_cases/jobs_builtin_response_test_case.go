package test_cases

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/screen_state"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// jobsBuiltinOutputLineRegex matches a single jobs output line and captures:
// 1. job id (integer), 2. marker (+ or - or space), 3. status (non-whitespace), 4. launch command (rest of line)
var jobsBuiltinOutputLineRegex = regexp.MustCompile(`^\[(\d+)\]([\+\-\s])\s+(\S+)\s+(.+)$`)

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

// validateJobsOutputLine checks one line of jobs output: matches the format regex, then validates
// captured job id, marker, status and launch command against the expected entry.
func (t JobsBuiltinResponseTestCase) validateJobsOutputLine(screenState screen_state.ScreenState, startRowIndex int, expectedEntry JobsBuiltinOutputEntry) (processedRowCount int, err *assertions.AssertionError) {
	processedRowCount = 1
	row := screenState.GetRow(startRowIndex)
	line := row.String()

	submatches := jobsBuiltinOutputLineRegex.FindStringSubmatch(line)
	if submatches == nil {
		message := "Line does not match expected jobs output format (expected: [job_id](+|-| )  status  launch_command)."
		if startRowIndex > screenState.GetLastLoggableRowIndex() {
			message = "Didn't find expected line.\n" + message
		} else {
			message = "Line does not match expected pattern.\n" + message + "\nGot: " + line
		}
		return 0, &assertions.AssertionError{ErrorRowIndex: startRowIndex, Message: message}
	}

	capturedJobIDStr := submatches[1]
	capturedMarkerStr := submatches[2]
	capturedStatusStr := submatches[3]
	capturedLaunchCommandStr := submatches[4]

	capturedJobID, parseErr := strconv.Atoi(capturedJobIDStr)
	if parseErr != nil {
		return 0, &assertions.AssertionError{
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Job id is not an integer: %q.", capturedJobIDStr),
		}
	}

	capturedMarker := parsedMarkerToMarkerConstant(capturedMarkerStr)

	// Normalise launch command: expected entry has no trailing " &"; shell output may include " &"
	normalisedCapturedLaunchCommand := strings.TrimSuffix(strings.TrimSpace(capturedLaunchCommandStr), "&")
	normalisedCapturedLaunchCommand = strings.TrimSpace(normalisedCapturedLaunchCommand)
	expectedLaunchCommandNormalised := strings.TrimSpace(expectedEntry.LaunchCommand)

	if capturedJobID != expectedEntry.JobNumber {
		return 0, &assertions.AssertionError{
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Job number mismatch: expected %d, got %d.", expectedEntry.JobNumber, capturedJobID),
		}
	}
	if capturedMarker != expectedEntry.Marker {
		return 0, &assertions.AssertionError{
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Marker mismatch: expected %s, got %s.",
				markerConstantToDisplay(expectedEntry.Marker), markerConstantToDisplay(capturedMarker)),
		}
	}
	if !strings.EqualFold(capturedStatusStr, expectedEntry.Status) {
		return 0, &assertions.AssertionError{
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Status mismatch: expected %q, got %q.", expectedEntry.Status, capturedStatusStr),
		}
	}
	if normalisedCapturedLaunchCommand != expectedLaunchCommandNormalised {
		return 0, &assertions.AssertionError{
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Launch command mismatch: expected %q, got %q.",
				expectedEntry.LaunchCommand, capturedLaunchCommandStr),
		}
	}

	return processedRowCount, nil
}

// jobsBuiltinOutputLinesAssertion runs validateJobsOutputLine for each expected entry.
type jobsBuiltinOutputLinesAssertion struct {
	tc *JobsBuiltinResponseTestCase
}

func (a jobsBuiltinOutputLinesAssertion) Inspect() string {
	return fmt.Sprintf("JobsBuiltinOutputLinesAssertion (%d lines)", len(a.tc.ExpectedOutputItems))
}

func (a jobsBuiltinOutputLinesAssertion) Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *assertions.AssertionError) {
	for _, expectedEntry := range a.tc.ExpectedOutputItems {
		n, assertionErr := a.tc.validateJobsOutputLine(screenState, startRowIndex+processedRowCount, expectedEntry)
		if assertionErr != nil {
			return processedRowCount, assertionErr
		}
		processedRowCount += n
	}
	return processedRowCount, nil
}

// parsedMarkerToMarkerConstant converts a single rune string "+", "-" or " " to UnmarkedJob, CurrentJob, PreviousJob.
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

// markerConstantToDisplay returns a display string for the marker constant for error messages.
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

func (t JobsBuiltinResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	command := "jobs"

	if err := shell.SendCommand(command); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", command)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	asserter.AddAssertion(jobsBuiltinOutputLinesAssertion{tc: &t})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
