package assertions

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/utils"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
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

func (a SingleLineAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	if a.ExpectedOutput == "" {
		panic("CodeCrafters Internal Error: ExpectedOutput must be provided")
	}

	processedRowCount = 1
	if a.StayOnSameLine {
		processedRowCount = 0
	}

	rawRow := screenState[startRowIndex]
	cleanedRow := virtual_terminal.BuildCleanedRow(rawRow)
	for _, pattern := range a.FallbackPatterns {
		if pattern.Match([]byte(cleanedRow)) {
			return processedRowCount, nil
		}
	}

	if cleanedRow != a.ExpectedOutput {
		detailedErrorMessage := utils.BuildColoredErrorMessage(a.ExpectedOutput, cleanedRow)
		return 0, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       "Output does not match expected value.\n" + detailedErrorMessage,
		}
	} else {
		return processedRowCount, nil
	}
}
