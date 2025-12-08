package assertions

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/screen_state"
)

// MultiLineAssertion asserts that multiple lines of output matches against a given array of strings
// Or a multi-line regex pattern(s)
type MultiLineAssertion struct {
	SingleLineAssertions []SingleLineAssertion
}

func NewMultiLineAssertion(expectedOutput []string) MultiLineAssertion {
	// No way to add fallbackPatterns through this constructor
	singleLineAssertions := []SingleLineAssertion{}
	for _, expectedLine := range expectedOutput {
		singleLineAssertions = append(singleLineAssertions, SingleLineAssertion{
			ExpectedOutput: expectedLine,
		})
	}

	return MultiLineAssertion{
		SingleLineAssertions: singleLineAssertions,
	}
}

func NewEmptyMultiLineAssertion() MultiLineAssertion {
	return MultiLineAssertion{
		SingleLineAssertions: []SingleLineAssertion{},
	}
}

// AddSingleLineAssertion is the recommended way to add single line assertions
// When they contain fallbackPatterns
func (a *MultiLineAssertion) AddSingleLineAssertion(expectedOutput string, fallbackPatterns []*regexp.Regexp) *MultiLineAssertion {
	a.SingleLineAssertions = append(a.SingleLineAssertions, SingleLineAssertion{
		ExpectedOutput:   expectedOutput,
		FallbackPatterns: fallbackPatterns,
	})
	return a
}

func (a *MultiLineAssertion) Inspect() string {
	return fmt.Sprintf("MultiLineAssertion (%v)", a.SingleLineAssertions)
}

func (a *MultiLineAssertion) Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *AssertionError) {
	totalProcessedRowCount := 0

	for i := range a.SingleLineAssertions {
		processedRowCount, err = a.SingleLineAssertions[i].Run(screenState, startRowIndex+totalProcessedRowCount)
		if err != nil {
			return totalProcessedRowCount, err
		}
		totalProcessedRowCount += processedRowCount
	}
	return totalProcessedRowCount, nil
}
