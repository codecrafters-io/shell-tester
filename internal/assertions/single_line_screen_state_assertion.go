package assertions

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/utils"
)

// SingleLineScreenStateAssertion are implicitly constrained to a single line of output
// Our ScreenState is composed of multiple lines, so we need to assert on each line individually
// This SingleLineScreenStateAssertion will assert only on a single given row (rowIndex)
// Ideally, we want to be able to assert using the expectedOutput string, a == b matching
// But, if that is not possible, we can use fallbackPatterns to match against multiple regexes
// And in the failure case, we want to show the expectedPatternExplanation to the user
type SingleLineScreenStateAssertion struct {
	// expectedOutput is the expected output string to match against
	expectedOutput string

	// fallbackPatterns is a list of regex patterns to match against
	fallbackPatterns []*regexp.Regexp

	// expectedPatternExplanation is the explanation of the expected pattern to
	// show in the error message in case of failure
	expectedPatternExplanation string
}

func NewSingleLineScreenStateAssertion(screenAsserter *ScreenAsserter, rowIndex int, expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string) SingleLineScreenStateAssertion {
	return SingleLineScreenStateAssertion{expectedOutput: expectedOutput, fallbackPatterns: fallbackPatterns, expectedPatternExplanation: expectedPatternExplanation}
}

// ToDo: screenState as its own type and wrap index / cursors inside it
func (t SingleLineScreenStateAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err error) {
	screen := screenState
	processedRowCount = 1

	if len(screen) == 0 {
		return processedRowCount, fmt.Errorf("expected screen to have at least one row, but it was empty")
	}
	rawRow := screen[startRowIndex]
	cleanedRow := utils.BuildCleanedRow(rawRow)

	if t.fallbackPatterns != nil && t.expectedPatternExplanation == "" {
		// expectedPatternExplanation is required for the error message on the FallbackPatterns path
		panic("CodeCrafters Internal Error: expectedPatternExplanation is empty on FallbackPatterns path")
	}

	regexPatternMatch := false

	// For each fallback pattern, check if the output matches
	// If it does, we break out of the loop and don't check for anything else, just return nil
	for _, pattern := range t.fallbackPatterns {
		if pattern.Match([]byte(cleanedRow)) {
			regexPatternMatch = true
			break
		}
	}

	if !regexPatternMatch {
		// No regex match till now, if expectedOutput is nil, we need to return an error
		// On this path, expectedPatternExplanation is required for the error message
		if t.expectedOutput == "" {
			// ToDo: Can't log it here
			// As this assertion would repeatedly fail while reading bytes
			// Possibly change loggers / return from here log outside
			// detailedErrorMessage := BuildColoredErrorMessage(t.expectedPatternExplanation, cleanedRow)
			// t.screenAsserter.Logger.Infof(detailedErrorMessage)
			return processedRowCount, fmt.Errorf("Received output does not match expectation.")
		} else {
			// ExpectedOutput is not nil, we can use it for exact string comparison
			if cleanedRow != t.expectedOutput {
				// detailedErrorMessage := BuildColoredErrorMessage(t.expectedOutput, cleanedRow)
				// t.screenAsserter.Logger.Infof(detailedErrorMessage)
				return processedRowCount, fmt.Errorf("Received output does not match expectation.")
			}
		}
	}

	return processedRowCount, nil
}
