package assertions

import (
	"time"
)

const defaultBellTimeout = 100 * time.Millisecond

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

func (a BellAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	if checkIfBellReceived(a.BellChannel) {
		return 0, nil
	}
	return 0, &AssertionError{
		StartRowIndex: startRowIndex,
		ErrorRowIndex: startRowIndex,
		Message:       "Expected bell to ring, but it didn't",
	}
}

// checkIfBellReceived waits for a bell signal on the provided channel with a timeout.
// Returns true if a signal is received, false if the timeout is reached.
func checkIfBellReceived(bellChannel chan bool) bool {
	select {
	case <-bellChannel:
		return true
	case <-time.After(defaultBellTimeout):
		// This timeout is for reading the bell signal on the bellChannel
		// Reading from term & processing happens in x/vt
		return false
	}
}
