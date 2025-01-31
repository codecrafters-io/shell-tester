package assertions

import (
	"time"
)

// BellAssertion asserts that ...
type BellAssertion struct {
	// Vt *virtual_terminal.VirtualTerminal
	BellChannel chan bool
}

func (a BellAssertion) Inspect() string {
	return "BellAssertion"
}

func (a BellAssertion) Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError) {
	if !checkIfBellReceived(a.BellChannel) {
		return 0, &AssertionError{
			StartRowIndex: startRowIndex,
			ErrorRowIndex: startRowIndex,
			Message:       "Expected bell to ring, but it didn't",
		}
	} else {
		return 0, nil
	}
}

func checkIfBellReceived(bellChannel chan bool) bool {
	select {
	case <-bellChannel:
		return true
	case <-time.After(100 * time.Millisecond): // Add reasonable timeout
		return false
	}
}
