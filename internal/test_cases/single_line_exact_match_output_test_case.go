package test_cases

import (
	"fmt"
	"regexp"
	"unicode"

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
	FallbackPatterns []*regexp.Regexp

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

func BuildExactMatchValidator(fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string, expectedOutput string, logger *logger.Logger) func([]byte) error {
	return func(output []byte) error {
		if fallbackPatterns != nil && expectedPatternExplanation == "" {
			// expectedPatternExplanation is required for the error message on the FallbackPatterns path
			panic("CodeCrafters Internal Error: expectedPatternExplanation is empty on FallbackPatterns path")
		}

		regexPatternMatch := false

		// For each fallback pattern, check if the output matches
		// If it does, we break out of the loop and don't check for anything else, just return nil
		for _, pattern := range fallbackPatterns {
			if pattern.Match(output) {
				regexPatternMatch = true
				break
			}
		}

		if !regexPatternMatch {
			// No regex match till now, if expectedOutput is nil, we need to return an error
			// On this path, expectedPatternExplanation is required for the error message
			if expectedOutput == "" {
				detailedErrorMessage := BuildColoredErrorMessage(expectedPatternExplanation, string(output))
				logger.Infof(detailedErrorMessage)
				return fmt.Errorf("Received output does not match expectation.")
			} else {
				// ExpectedOutput is not nil, we can use it for exact string comparison
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

func BuildColoredErrorMessage(expectedPatternExplanation string, output string) string {
	errorMsg := colorizeString(color.FgGreen, "Expected:")
	errorMsg += " \"" + expectedPatternExplanation + "\""
	errorMsg += "\n"
	errorMsg += colorizeString(color.FgRed, "Received:")
	errorMsg += " \"" + removeNonPrintableCharacters(output) + "\""

	return errorMsg
}

func removeNonPrintableCharacters(output string) string {
	result := ""
	for _, r := range output {
		if unicode.IsPrint(r) {
			result += string(r)
		} else {
			result += "�" // U+FFFD
		}
	}
	return result
}
