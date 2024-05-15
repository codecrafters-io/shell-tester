package assertions

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
)

type BufferAssertion struct {
	ExpectedValue string
	ActualValue   string

	stdOutOffset                      int
	stderrOffset                      int
	previousOperationDataPipeSelector string
}

func (a *BufferAssertion) Run(executable *shell_executable.ShellExecutable, dataPipeSelector string) error {
	bytesValue, err := executable.ReadBuffer(dataPipeSelector)
	if err != nil {
		return err
	}

	var idx int
	switch dataPipeSelector {
	case "stdout":
		idx = a.stdOutOffset
	case "stderr":
		idx = a.stderrOffset
	}

	currentOperationBytesValue := bytesValue[idx:]
	value := string(removeControlSequence(currentOperationBytesValue))
	a.previousOperationDataPipeSelector = dataPipeSelector
	a.ActualValue = value

	if len(value) == 0 {
		return fmt.Errorf("Expected to receive value, but got nothing")
	}

	if !strings.Contains(value, a.ExpectedValue) {
		return fmt.Errorf("Expected value to be %q, but got %q", a.ExpectedValue, value)
	}
	return nil
}

func (a *BufferAssertion) UpdateOffsetToCurrentLength() {
	if a.previousOperationDataPipeSelector == "stdout" {
		a.stdOutOffset += len(a.ActualValue)
	}
	if a.previousOperationDataPipeSelector == "stderr" {
		a.stderrOffset += len(a.ActualValue)
	}
}

// Doesn't fit nicely with BufferAssertion
// type TruncatedBuffer struct {
// 	buffer bytes.Buffer
// 	offset int
// }
// func (t TruncatedBuffer) updateOffsetToCurrentLength() {
// 	t.offset = t.buffer.Len()
// }

// Duplicated for now
func removeControlSequence(data []byte) []byte {
	PROMPT_START := '$'

	for startIdx, r := range string(data) {
		// Standard escape codes are prefixed with Escape (27)
		if r == 27 {
			// remove from here upto PROMPT_START
			for endIdx, r2 := range string(data[startIdx:]) {
				if r2 == PROMPT_START {
					// Remove from start_idx to end_idx-1
					data = append(data[:startIdx], data[endIdx:]...)
					break
				}
			}
		}
	}

	return data
}
