package utils

import (
	"unicode"

	"github.com/fatih/color"
)

const VT_SENTINEL_CHARACTER = "."

func ColorizeString(colorToUse color.Attribute, msg string) string {
	c := color.New(colorToUse)
	return c.Sprint(msg)
}

func BuildColoredErrorMessage(expectedPatternExplanation string, output string) string {
	errorMsg := ColorizeString(color.FgGreen, "Expected:")
	errorMsg += " \"" + expectedPatternExplanation + "\""
	errorMsg += "\n"
	errorMsg += ColorizeString(color.FgRed, "Received:")
	errorMsg += " \"" + RemoveNonPrintableCharacters(output) + "\""

	return errorMsg
}

func RemoveNonPrintableCharacters(output string) string {
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

// TODO: What if there are tabs, will we have SENTINEL_CHAR in the middle of a row?
// Alternative here would be to TrimRight(SENTINEL_CHAR)
func BuildCleanedRow(row []string) string {
	result := ""
	for _, cell := range row {
		if cell != VT_SENTINEL_CHARACTER {
			result += cell
		}
	}
	return result
}

func AsBool(T func() error) func() bool {
	// Takes in a function that takes no params & returns an error
	// Returns the function wrapped in a helper such that it returns a bool
	// in liue of the error, true if the function execution is a success
	return func() bool { return T() == nil }
}
