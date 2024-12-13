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

func (a SingleLineAssertion) Inspect() string {
	return fmt.Sprintf("SingleLineAssertion (%q)", a.ExpectedOutput)
}

func (a SingleLineAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
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
		// TODO: Review
		// detailedErrorMessage := utils.BuildColoredErrorMessage(a.ExpectedOutput, cleanedRow)
		return 0, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       fmt.Sprintf("Expected %q, got %q", a.ExpectedOutput, cleanedRow),
			// Message:       detailedErrorMessage + "\n" + fmt.Sprintf("Expected %q, got %q", a.ExpectedOutput, cleanedRow),
		}
	} else {
		return 1, nil
	}
}
