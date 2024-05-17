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

func (a *BufferAssertion) Run(buffer *shell_executable.TruncatedBuffer, coreTest func(string, []byte) error) error {
	bytesValue, err := buffer.ReadBuffer(coreTest, a.ExpectedValue)
	if err != nil {
		return err
	}

	value := string(shell_executable.RemoveControlSequence(bytesValue))
	value = strings.TrimSpace(value)
	a.ActualValue = value

	return coreTest(a.ExpectedValue, bytesValue)
}

func CoreTestInexact(expectedData string, actualData []byte) error {
	value := string(shell_executable.RemoveControlSequence(actualData))
	value = strings.TrimSpace(value)

	if len(value) == 0 {
		return fmt.Errorf("Expected to receive value, but got nothing")
	}

	if !strings.Contains(value, expectedData) {
		// ToDo: Update log accordingly "contains"
		return fmt.Errorf("Expected value to be %q, but got %q", expectedData, value)
	}

	return nil
}

func CoreTestExact(expectedData string, actualData []byte) error {
	value := string(shell_executable.RemoveControlSequence(actualData))
	value = strings.TrimSpace(value)

	if len(value) == 0 {
		return fmt.Errorf("Expected to receive value, but got nothing")
	}

	if !strings.EqualFold(value, expectedData) {
		return fmt.Errorf("Expected value to be %q, but got %q", expectedData, value)
	}

	return nil
}
