package assertions

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
)

type BufferAssertion struct {
	ExpectedValue string
	ActualValue   string
}

func (a *BufferAssertion) Run(buffer *shell_executable.TruncatedBuffer) error {
	bytesValue, err := buffer.ReadBuffer()
	if err != nil {
		return err
	}

	value := string(shell_executable.RemoveControlSequence(bytesValue))
	a.ActualValue = value

	if len(value) == 0 {
		return fmt.Errorf("Expected to receive value, but got nothing")
	}

	if !strings.Contains(value, a.ExpectedValue) {
		// ToDo: Update log accordingly "contains"
		return fmt.Errorf("Expected value to be %q, but got %q", a.ExpectedValue, value)
	}
	return nil
}
