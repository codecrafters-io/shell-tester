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
	ExpectedOutput string

	// FallbackPatterns are a list of regex patterns that are evaluated against the command's output first
	// and only if none of them match, the ExpectedOutput is used to compare against the output.
	FallbackPatterns []regexp.Regexp

	// ExpectedPatternExplanation is used in the error message if the ExpectedPattern doesn't match the command's output
	ExpectedPatternExplanation string

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t SingleLineExactMatchTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	singleLineOutputTestCase := singleLineOutputTestCase{
		Command:        t.Command,
		Validator:      BuildExactMatchValidator(t.FallbackPatterns, t.ExpectedPatternExplanation, t.ExpectedOutput, logger),
		SuccessMessage: t.SuccessMessage,
	}

	return singleLineOutputTestCase.Run(shell, logger)

}

func BuildExactMatchValidator(fallbackPatterns []regexp.Regexp, expectedOutputExplanation string, expectedOutput string, logger *logger.Logger) func([]byte) error {
	return func(output []byte) error {
		regexPatternMatch := false

		for _, pattern := range fallbackPatterns {
			if pattern.Match(output) {
				regexPatternMatch = true
				break
			}
		}

		if !regexPatternMatch {
			if expectedOutput == "" {
				detailedErrorMessage := BuildColoredErrorMessage(expectedOutputExplanation, string(output))
				logger.Infof(detailedErrorMessage)
				return fmt.Errorf("Received output does not match expectation.")
			} else {
				if string(output) != expectedOutput {
					detailedErrorMessage := BuildColoredErrorMessage(expectedOutput, string(output))
					logger.Infof(detailedErrorMessage)
					return fmt.Errorf("Received output does not match expectation.")
				}
			}
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
	errorMsg += " \"" + expectedPatternExplanation + "\""
	errorMsg += "\n"
	errorMsg += colorizeString(color.FgRed, "Received:")
	errorMsg += " \"" + cleanedOutput + "\""

	return errorMsg
}
