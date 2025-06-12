package assertions

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/screen_state"
	"github.com/codecrafters-io/shell-tester/internal/utils"
)

// SingleLineAssertion asserts that a single line of output matches a given string or regex pattern(s)
type SingleLineAssertion struct {
	// ExpectedOutput is the expected output string to match against
	ExpectedOutput string

	// FallbackPatterns is a list of regex patterns to match against. This is useful to handle shell-specific variable behaviour
	FallbackPatterns []*regexp.Regexp

	// StayOnSameLine is a flag to indicate that the shell cursor
	// should stay on the same line after the assertion is run
	// Most probably because the next assertion will run on the same line
	StayOnSameLine bool
}

func (a SingleLineAssertion) Inspect() string {
	return fmt.Sprintf("SingleLineAssertion (%q)", a.ExpectedOutput)
}

func (a SingleLineAssertion) Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *AssertionError) {
	if a.ExpectedOutput == "" && len(a.FallbackPatterns) == 0 {
		panic("CodeCrafters Internal Error: ExpectedOutput or fallbackPatterns must be provided")
	}

	processedRowCount = 1
	if a.StayOnSameLine {
		processedRowCount = 0
	}

	row := screenState.GetRow(startRowIndex)

	for _, pattern := range a.FallbackPatterns {
		if pattern.Match([]byte(row.String())) {
			return processedRowCount, nil
		}
	}

	if row.String() != a.ExpectedOutput {
		rowDescription := ""

		if startRowIndex > screenState.GetLastLoggableRowIndex() {
			rowDescription = "no line received"
		} else if row.IsEmpty() {
			rowDescription = "empty line"
		}

		detailedErrorMessage := utils.BuildColoredErrorMessage(a.ExpectedOutput, row.String(), rowDescription)

		// If the line won't be logged, we say "didn't find line ..." instead of "line does not match expected ..."
		if startRowIndex > screenState.GetLastLoggableRowIndex() {
			return 0, &AssertionError{
				ErrorRowIndex: startRowIndex,
				Message:       "Didn't find expected line.\n" + detailedErrorMessage,
			}
		} else {
			return 0, &AssertionError{
				ErrorRowIndex: startRowIndex,
				Message:       "Line does not match expected value.\n" + detailedErrorMessage,
			}
		}
	} else {
		return processedRowCount, nil
	}
}
