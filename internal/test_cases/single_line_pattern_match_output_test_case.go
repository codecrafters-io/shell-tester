package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// SingleLinePatternMatchTestCase verifies a prompt exists, sends a command and matches the output against a string.
type SingleLinePatternMatchTestCase struct {
	// The command to execute (the command's output will be matched against ExpectedPattern)
	Command string

	// ExpectedPattern is the regex that is evaluated against the command's output.
	ExpectedPattern string

	// ExpectedPatternExplanation is used in the error message if the ExpectedPattern doesn't match the command's output
	ExpectedPatternExplanation string

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t SingleLinePatternMatchTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	singleLineOutputTestCase := SingleLineOutputTestCase{
		Command:        t.Command,
		Validator:      BuildPatternMatchValidator(t.ExpectedPattern, t.ExpectedPatternExplanation),
		SuccessMessage: t.SuccessMessage,
	}

	return singleLineOutputTestCase.Run(shell, logger)

}

func BuildPatternMatchValidator(pattern string, simplifiedPatternExplanation string) func([]byte) error {
	re := regexp.MustCompile(pattern)
	return func(output []byte) error {
		if !re.Match(output) {
			return fmt.Errorf("Expected first line of output to contain %s, got %q", simplifiedPatternExplanation, string(output))
		}
		return nil
	}
}
