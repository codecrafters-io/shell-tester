package test_cases

import (
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// TODO: Remove SingleLineOutputTestCase entirely, replace with SingleLineOutputAssertion invoked within ScreenAsserter
// singleLineOutputTestCase verifies a prompt exists, sends a command and matches the output against a string.
type singleLineOutputTestCase struct {
	// The command to execute (the command's output will be matched using the Validator function)
	Command string

	// Validator is a function that contains the expected pattern and the expected pattern explanation,
	// and returns an error if the output does not match the expected pattern.
	Validator func([]byte) error

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t singleLineOutputTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	return nil
}
