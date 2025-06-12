package assertions

import "github.com/codecrafters-io/shell-tester/internal/screen_state"

type Assertion interface {
	Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *AssertionError)
	Inspect() string
}
