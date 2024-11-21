package test_cases

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/fatih/color"
)

// SingleLineExactMatchTestCase verifies a prompt exists, sends a command and matches the output against a string.
type SingleLineExactMatchTestCase struct {
	// The command to execute (the command's output will be matched against ExpectedPattern)
	Command string

	// ExpectedPattern is the regex that is evaluated against the command's output.
	ExpectedPattern string

	// ExpectedPatternExplanation is used in the error message if the ExpectedPattern doesn't match the command's output
	ExpectedPatternExplanation string

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t SingleLineExactMatchTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	singleLineOutputTestCase := SingleLineOutputTestCase{
		Command:        t.Command,
		Validator:      BuildExactMatchValidator(t.ExpectedPattern, t.ExpectedPatternExplanation),
		SuccessMessage: t.SuccessMessage,
	}

	return singleLineOutputTestCase.Run(shell, logger)

}

func BuildExactMatchValidator(pattern string, simplifiedPatternExplanation string) func([]byte) error {
	re := regexp.MustCompile(pattern)
	return func(output []byte) error {
		if !re.Match(output) {
			return fmt.Errorf(BuildColoredErrorMessage(simplifiedPatternExplanation, string(output)))
		}
		return nil
	}
}

func colorizeString(colorToUse color.Attribute, msg string) string {
	c := color.New(colorToUse)
	return c.Sprint(msg)
}

func BuildColoredErrorMessage(expectedPatternExplanation string, cleanedOutput string) string {
	indent := 9 - 4

	errorMsg := "Expected:" // 9
	errorMsg += " " + colorizeString(color.FgGreen, expectedPatternExplanation)
	errorMsg += "\n"
	errorMsg += strings.Repeat(" ", indent)
	errorMsg += "Got:" // 4
	errorMsg += " " + colorizeString(color.FgRed, cleanedOutput)

	return errorMsg
}
