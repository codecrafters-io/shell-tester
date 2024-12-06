package assertions

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/fatih/color"
)

type ScreenStateAssertion struct {
	expectedOutput             string
	fallbackPatterns           []*regexp.Regexp
	expectedPatternExplanation string
}

func NewScreenStateAssertion(expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string) ScreenStateAssertion {
	return ScreenStateAssertion{expectedOutput: expectedOutput, fallbackPatterns: fallbackPatterns, expectedPatternExplanation: expectedPatternExplanation}
}

func (a ScreenStateAssertion) Run(output string, logger *logger.Logger) error {
	if a.fallbackPatterns != nil && a.expectedPatternExplanation == "" {
		// expectedPatternExplanation is required for the error message on the FallbackPatterns path
		panic("CodeCrafters Internal Error: expectedPatternExplanation is empty on FallbackPatterns path")
	}

	regexPatternMatch := false

	// For each fallback pattern, check if the output matches
	// If it does, we break out of the loop and don't check for anything else, just return nil
	for _, pattern := range a.fallbackPatterns {
		if pattern.Match([]byte(output)) {
			regexPatternMatch = true
			break
		}
	}

	if !regexPatternMatch {
		// No regex match till now, if expectedOutput is nil, we need to return an error
		// On this path, expectedPatternExplanation is required for the error message
		if a.expectedOutput == "" {
			detailedErrorMessage := BuildColoredErrorMessage(a.expectedPatternExplanation, output)
			logger.Infof(detailedErrorMessage)
			return fmt.Errorf("Received output does not match expectation.")
		} else {
			// ExpectedOutput is not nil, we can use it for exact string comparison
			if output != a.expectedOutput {
				detailedErrorMessage := BuildColoredErrorMessage(a.expectedOutput, output)
				logger.Infof(detailedErrorMessage)
				return fmt.Errorf("Received output does not match expectation.")
			}
		}
	}

	return nil
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
