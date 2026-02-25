package assertions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/screen_state"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/dustin/go-humanize/english"
)

// SingleLineRegexAssertion asserts that a single line of output matches a given set of regexes
type SingleLineRegexAssertion struct {
	ExpectedRegexPatterns []*regexp.Regexp

	// StayOnSameLine is a flag to indicate that the shell cursor
	// should stay on the same line after the assertion is run
	// Most probably because the next assertion will run on the same line
	StayOnSameLine bool
}

func (a SingleLineRegexAssertion) Inspect() string {
	patterns := []string{}

	for _, pattern := range a.ExpectedRegexPatterns {
		patterns = append(patterns, pattern.String())
	}

	return fmt.Sprintf("SingleLineAssertion (%q)", strings.Join(patterns, "\n"))
}

func (a SingleLineRegexAssertion) Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *AssertionError) {
	if len(a.ExpectedRegexPatterns) == 0 {
		panic("CodeCrafters Internal Error: ExpectedRegexPatterns must be provided")
	}

	processedRowCount = 1
	if a.StayOnSameLine {
		processedRowCount = 0
	}

	row := screenState.GetRow(startRowIndex)

	for _, pattern := range a.ExpectedRegexPatterns {
		if pattern.Match([]byte(row.String())) {
			return processedRowCount, nil
		}
	}

	var rowDescription string

	if startRowIndex > screenState.GetLastLoggableRowIndex() {
		rowDescription = "no line received"
	} else if row.IsEmpty() {
		rowDescription = "empty line"
	}

	detailedErrorMessage := utils.BuildColoredErrorMessageForFallbackPatternMismatch(a.ExpectedRegexPatterns, row.String(), rowDescription)

	// Line does not exist at all
	if startRowIndex > screenState.GetLastLoggableRowIndex() {
		return 0, &AssertionError{
			ErrorRowIndex: startRowIndex,
			Message:       "Didn't find expected line.\n" + detailedErrorMessage,
		}
	}

	// Line exists: but doesn't match any of the expected regexes
	return 0, &AssertionError{
		ErrorRowIndex: startRowIndex,
		Message: fmt.Sprintf("Line does not match %s.\n%s",
			english.PluralWord(len(a.ExpectedRegexPatterns), "the expected regex", "any of the expected regexes"),
			detailedErrorMessage,
		),
	}
}
