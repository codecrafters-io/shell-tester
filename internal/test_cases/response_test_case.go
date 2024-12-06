package test_cases

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// TODO: Remove ResponseTestCase entirely, replace with SingleLineOutputAssertion invoked within ScreenAsserter
// ResponseTestCase reads the output from the shell, and verifies that it matches the expected output.
type ResponseTestCase struct {
	// expectedPrompt is the prompt expected to be displayed (example: "$ ")
	expectedPrompt string

	assertion assertions.SingleLineScreenStateAssertion

	shouldOmitSuccessLog bool
}

func NewResponseTestCase(expectedPrompt string, assertion assertions.SingleLineScreenStateAssertion, shouldOmitSuccessLog bool) ResponseTestCase {
	return ResponseTestCase{expectedPrompt: expectedPrompt, assertion: assertion, shouldOmitSuccessLog: shouldOmitSuccessLog}
}

func NewSilentResponseTestCase(expectedPrompt string, assertion assertions.SingleLineScreenStateAssertion) ResponseTestCase {
	return ResponseTestCase{expectedPrompt: expectedPrompt, assertion: assertion, shouldOmitSuccessLog: true}
}

func (t ResponseTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	return nil
}
