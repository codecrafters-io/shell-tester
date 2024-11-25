package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/fatih/color"
)

// SingleLineExactMatchTestCase internally creates a SingleLineOutputTestCase with a validator that matches the output against a pattern.
// We look for an exact match in this case, so our error logs can take advantage of color to highlight the exact mismatch.
type SingleLineExactMatchTestCase struct {
	// Command is the command to execute, whose output will be matched against ExpectedPattern.
	Command string

	// ExpectedPattern is the regex that is evaluated against the command's output.
	ExpectedPattern string

	// ExpectedPatternExplanation is used in the error message if the ExpectedPattern doesn't match the command's output
	ExpectedPatternExplanation string

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t SingleLineExactMatchTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	singleLineOutputTestCase := singleLineOutputTestCase{
		Command:        t.Command,
		Validator:      BuildExactMatchValidator(t.ExpectedPattern, t.ExpectedPatternExplanation, logger),
		SuccessMessage: t.SuccessMessage,
	}

	return singleLineOutputTestCase.Run(shell, logger)

}

func BuildExactMatchValidator(pattern string, simplifiedPatternExplanation string, logger *logger.Logger) func([]byte) error {
	re := regexp.MustCompile(pattern)
	return func(output []byte) error {
		if !re.Match(output) {
			detailedErrorMessage := BuildColoredErrorMessage(simplifiedPatternExplanation, string(output))
			logger.Infof(detailedErrorMessage)
			return fmt.Errorf("Received output does not match expectation.")
		}
		return nil
	}
}

func colorizeString(colorToUse color.Attribute, msg string) string {
	c := color.New(colorToUse)
	return c.Sprint(msg)
}

func BuildColoredErrorMessage(expectedPatternExplanation string, cleanedOutput string) string {
	errorMsg := colorizeString(color.FgGreen, "Expected:")
	errorMsg += fmt.Sprintf("%q", expectedPatternExplanation)
	errorMsg += "\n"
	errorMsg += colorizeString(color.FgRed, "Received:")
	errorMsg += fmt.Sprintf("%q", cleanedOutput)

	return errorMsg
}
