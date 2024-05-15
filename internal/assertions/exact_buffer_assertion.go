package assertions

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
)

type ExactBufferAssertion struct {
	ExpectedValue string
	ActualValue   string
}

func (a *ExactBufferAssertion) Run(executable *shell_executable.ShellExecutable, dataPipeSelector string) error {
	bytesValue, err := executable.ReadBuffer(dataPipeSelector)
	if err != nil {
		return err
	}

	value := string(removeControlSequence(bytesValue))
	a.ActualValue = value

	if len(value) == 0 {
		return fmt.Errorf("Expected to receive value, but got nothing")
	}

	if !strings.EqualFold(value, a.ExpectedValue) {
		return fmt.Errorf("Expected value to be %q, but got %q", a.ExpectedValue, value)
	}

	return nil
}
