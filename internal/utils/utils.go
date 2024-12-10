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

func BuildCleanedRow(row []string) string {
	result := ""
	for _, cell := range row {
		if cell != VT_SENTINEL_CHARACTER {
			result += cell
		}
	}
	return result
}
