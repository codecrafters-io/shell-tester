package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// SingleLineStringMatchTestCase internally creates a SingleLineOutputTestCase with a validator that matches the output against a pattern.
// We look for an exact match in this case, so our error logs can take advantage of color to highlight the exact mismatch.
type SingleLineStringMatchTestCase struct {
	// Command is the command to execute, whose output will be matched against ExpectedPattern.
	Command string

	// ExpectedOutput is the regex that is evaluated against the command's output.
	ExpectedOutput string

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t SingleLineStringMatchTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	singleLineOutputTestCase := singleLineOutputTestCase{
		Command:        t.Command,
		Validator:      BuildExactStringMatchValidator(t.ExpectedOutput, logger),
		SuccessMessage: t.SuccessMessage,
	}

	return singleLineOutputTestCase.Run(shell, logger)

}

func BuildExactStringMatchValidator(expectedOutput string, logger *logger.Logger) func([]byte) error {
	return func(output []byte) error {
		if string(output) != expectedOutput {
			detailedErrorMessage := BuildColoredErrorMessage(expectedOutput, string(output))
			logger.Infof(detailedErrorMessage)
			return fmt.Errorf("Received output does not match expectation.")
		}
		return nil
	}
}
