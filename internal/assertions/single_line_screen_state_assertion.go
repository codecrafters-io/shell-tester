package assertions

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/fatih/color"
)

const VT_SENTINEL_CHARACTER = "."

// SingleLineScreenStateAssertion are implicitly constrained to a single line of output
// Our ScreenState is composed of multiple lines, so we need to assert on each line individually
// This SingleLineScreenStateAssertion will assert only on a single given row (rowIndex)
// Ideally, we want to be able to assert using the expectedOutput string, a == b matching
// But, if that is not possible, we can use fallbackPatterns to match against multiple regexes
// And in the failure case, we want to show the expectedPatternExplanation to the user
type SingleLineScreenStateAssertion struct {
	rowIndex                   int
	expectedOutput             string
	fallbackPatterns           []*regexp.Regexp
	expectedPatternExplanation string
}

func NewSingleLineScreenStateAssertion(rowIndex int, expectedOutput string, fallbackPatterns []*regexp.Regexp, expectedPatternExplanation string) SingleLineScreenStateAssertion {
	return SingleLineScreenStateAssertion{rowIndex: rowIndex, expectedOutput: expectedOutput, fallbackPatterns: fallbackPatterns, expectedPatternExplanation: expectedPatternExplanation}
}

// ToDo: screenState as its own type and wrap index / cursors inside it
func (a SingleLineScreenStateAssertion) Run(screenState [][]string, logger *logger.Logger) error {
	rawRow := screenState[a.rowIndex]
	cleanedRow := buildCleanedRow(rawRow)

	if a.fallbackPatterns != nil && a.expectedPatternExplanation == "" {
		// expectedPatternExplanation is required for the error message on the FallbackPatterns path
		panic("CodeCrafters Internal Error: expectedPatternExplanation is empty on FallbackPatterns path")
	}

	regexPatternMatch := false

	// For each fallback pattern, check if the output matches
	// If it does, we break out of the loop and don't check for anything else, just return nil
	for _, pattern := range a.fallbackPatterns {
		if pattern.Match([]byte(cleanedRow)) {
			regexPatternMatch = true
			break
		}
	}

	if !regexPatternMatch {
		// No regex match till now, if expectedOutput is nil, we need to return an error
		// On this path, expectedPatternExplanation is required for the error message
		if a.expectedOutput == "" {
			detailedErrorMessage := BuildColoredErrorMessage(a.expectedPatternExplanation, cleanedRow)
			logger.Infof(detailedErrorMessage)
			return fmt.Errorf("Received output does not match expectation.")
		} else {
			// ExpectedOutput is not nil, we can use it for exact string comparison
			if cleanedRow != a.expectedOutput {
				detailedErrorMessage := BuildColoredErrorMessage(a.expectedOutput, cleanedRow)
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

// ToDo: move this to its own package along with all vterm interface code
func buildCleanedRow(row []string) string {
	result := ""
	for _, cell := range row {
		if cell != VT_SENTINEL_CHARACTER {
			result += cell
		}
	}
	return result
}
