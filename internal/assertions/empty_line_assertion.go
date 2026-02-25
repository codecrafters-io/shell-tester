package assertions

import (
	"github.com/codecrafters-io/shell-tester/internal/screen_state"
	"github.com/codecrafters-io/shell-tester/internal/utils"
)

// EmptyLineAssertion asserts that a single line of output is empty
type EmptyLineAssertion struct {
	// StayOnSameLine is a flag to indicate that the shell cursor
	// should stay on the same line after the assertion is run
	// Most probably because the next assertion will run on the same line
	StayOnSameLine bool
}

func (a EmptyLineAssertion) Inspect() string {
	return "EmptyLineAssertion"
}

func (a EmptyLineAssertion) Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *AssertionError) {
	processedRowCount = 1
	if a.StayOnSameLine {
		processedRowCount = 0
	}

	row := screenState.GetRow(startRowIndex)

	if row.IsEmpty() {
		return processedRowCount, nil
	}

	// Build error message
	detailedErrorMessage := utils.BuildColoredErrorMessageForUnexpectedOutput("(empty line)", row.String(), "")
	message := "Line is not empty.\n" + detailedErrorMessage

	return 0, &AssertionError{
		ErrorRowIndex: startRowIndex,
		Message:       message,
	}
}
