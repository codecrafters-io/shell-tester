package assertions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/utils"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
)

// MultiLineAssertion asserts that multiple lines of output matches against a given array of strings
// Or a multi-line regex pattern(s)
type MultiLineAssertion struct {
	// ExpectedOutput is the array of expected output strings to match against
	ExpectedOutput []string

	// FallbackPatterns is a list of regex patterns to match against. This is useful to handle shell-specific variable behaviour
	FallbackPatterns []*regexp.Regexp
}

func (a MultiLineAssertion) Inspect() string {
	return fmt.Sprintf("MultiLineAssertion (%q)", a.ExpectedOutput)
}

func (a MultiLineAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	if len(a.ExpectedOutput) == 0 {
		panic("CodeCrafters Internal Error: ExpectedOutput must be provided")
	}

	totalRows := len(a.ExpectedOutput)
	rawRows := screenState[startRowIndex : startRowIndex+totalRows]
	cleanedRows := []string{}
	for _, rawRow := range rawRows {
		cleanedRows = append(cleanedRows, virtual_terminal.BuildCleanedRow(rawRow))
	}
	cleanedRowsString := strings.Join(cleanedRows, "\n")
	expectedOutputString := strings.Join(a.ExpectedOutput, "\n")

	for _, pattern := range a.FallbackPatterns {
		if pattern.Match([]byte(cleanedRowsString)) {
			return len(a.ExpectedOutput), nil
		}
	}

	if cleanedRowsString != expectedOutputString {
		detailedErrorMessage := utils.BuildColoredErrorMessage(expectedOutputString, cleanedRowsString)
		return 0, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       "Output does not match expected value.\n" + detailedErrorMessage,
		}
	} else {
		return len(a.ExpectedOutput), nil
	}
}
