package assertions

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/utils"
)

// SingleLineAssertion asserts that a single line of output matches a given string or regex pattern(s)
type SingleLineAssertion struct {
	// ExpectedOutput is the expected output string to match against
	ExpectedOutput string

	// FallbackPatterns is a list of regex patterns to match against. This is useful to handle shell-specific variable behaviour
	FallbackPatterns []*regexp.Regexp
}

func (a SingleLineAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	// TODO: Move these to assertion collection
	if len(screenState) == 0 {
		panic("CodeCrafters internal error: expected screen to have at least one row, but it was empty")
	}

	if startRowIndex >= len(screenState) {
		panic("CodeCrafters internal error: startRowIndex is larger than screenState rows")
	}

	if a.ExpectedOutput == "" {
		panic("CodeCrafters Internal Error: ExpectedOutput must be provided")
	}

	rawRow := screenState[startRowIndex]
	cleanedRow := utils.BuildCleanedRow(rawRow)

	for _, pattern := range a.FallbackPatterns {
		if pattern.Match([]byte(cleanedRow)) {
			return 1, nil
		}
	}

	if cleanedRow != a.ExpectedOutput {
		// TODO: Return colored error message when constructing AssertionError
		// detailedErrorMessage := BuildColoredErrorMessage(t.expectedPatternExplanation, cleanedRow)
		return 0, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Expected %q, got %q", a.ExpectedOutput, cleanedRow),
		}
	} else {
		return 1, nil
	}
}
