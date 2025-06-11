package assertions

import "github.com/codecrafters-io/shell-tester/internal/screen_state"

// BellAssertion asserts that the bell callback function is called
// by the virtual terminal, this can only happen if the user sends
// /U0007 to the shell.
type BellAssertion struct {
	// Vt *virtual_terminal.VirtualTerminal
	BellChannel chan bool
}

func (a BellAssertion) Inspect() string {
	return "BellAssertion"
}

func (a BellAssertion) Run(screenState screen_state.ScreenState, startRowIndex int) (processedRowCount int, err *AssertionError) {
	if checkIfBellReceived(a.BellChannel) {
		return 0, nil
	}
	return 0, &AssertionError{
		StartRowIndex: -1,
		ErrorRowIndex: -1,
		Message:       "Expected bell to ring, but it didn't",
	}
}

// checkIfBellReceived waits for a bell signal on the provided channel in a non-blocking way.
// Returns true if a signal is received, false if no signal is received or if the channel is closed.
func checkIfBellReceived(bellChannel chan bool) bool {
	select {
	case _, ok := <-bellChannel:
		if ok {
			return true
		}
		return false
	default:
		return false
	}
}
