package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/dustin/go-humanize/english"
	"github.com/fatih/color"
)

func ColorizeString(colorToUse color.Attribute, msg string) string {
	c := color.New(colorToUse)
	return c.Sprint(msg)
}

func BuildColoredErrorMessageForExpectedOutputMismatch(expectedOutput string, receivedOutput string, receivedOutputDescription string) string {
	errorMsg := ColorizeString(color.FgGreen, "Expected:")
	errorMsg += " \"" + expectedOutput + "\""
	errorMsg += "\n"
	errorMsg += ColorizeString(color.FgRed, "Received:")
	errorMsg += " \"" + RemoveNonPrintableCharacters(receivedOutput) + "\""

	if receivedOutputDescription != "" {
		errorMsg += " " + ColorizeString(color.FgRed, fmt.Sprintf("(%s)", receivedOutputDescription))
	}

	return errorMsg
}

func BuildColoredErrorMessageForFallbackPatternMismatch(fallbackPatterns []*regexp.Regexp, receivedOutput string, receivedOutputDescription string) string {
	if len(fallbackPatterns) == 0 {
		panic("Codecrafters Internal Error - BuildColoredErrorMessageForFallbackPatternMismatch called with empty regex array")
	}

	var errorMsg strings.Builder
	errorMsg.WriteString(ColorizeString(color.FgRed, "Received:"))
	errorMsg.WriteString(" \"" + RemoveNonPrintableCharacters(receivedOutput) + "\"")

	if receivedOutputDescription != "" {
		errorMsg.WriteString(" " + ColorizeString(color.FgRed, fmt.Sprintf("(%s)", receivedOutputDescription)))
	}

	errorMsg.WriteString("\n")

	expectedStatement := fmt.Sprintf("Expected line to match %s:\n",
		english.PluralWord(len(fallbackPatterns),
			"the following regex",
			"one of the following regexes"))

	errorMsg.WriteString(ColorizeString(color.FgGreen, expectedStatement))

	for _, regex := range fallbackPatterns {
		regexString := regex.String()
		hintString := fmt.Sprintf("Hint: %s", GetRegex101Link(regexString, receivedOutput))
		regexAndHint := fmt.Sprintf("%s (%s)", regexString, hintString)
		errorMsg.WriteString(ColorizeString(color.FgGreen, regexAndHint))
	}

	return errorMsg.String()
}

func GetRegex101Link(pattern string, testString string) string {
	return fmt.Sprintf("https://regex101.com/?regex=%s&testString=%s", url.QueryEscape(pattern), url.QueryEscape(testString))
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

func AsBool(T func() error) func() bool {
	// Takes in a function that takes no params & returns an error
	// Returns the function wrapped in a helper such that it returns a bool
	// in lieu of the error, true if the function execution is a success
	return func() bool { return T() == nil }
}

// LogReadableFileContents prints file contents in a readable way, replacing tabs and spaces with visible markers.
func LogReadableFileContents(l *logger.Logger, fileContents string, logMsg string, fileName string) {
	l.Infof("%s", logMsg)
	l.UpdateLastSecondaryPrefix(fileName)
	defer l.ResetSecondaryPrefixes()
	printableFileContents := strings.ReplaceAll(fileContents, "%", "%%")
	printableFileContents = strings.ReplaceAll(printableFileContents, "\t", "<|TAB|>")

	regex1 := regexp.MustCompile("[ ]+\n")
	regex2 := regexp.MustCompile("[ ]+$")
	printableFileContents = regex1.ReplaceAllString(printableFileContents, "\n")
	printableFileContents = regex2.ReplaceAllString(printableFileContents, "<|SPACE|>")

	if len(printableFileContents) == 0 {
		l.Plainf("<|EMPTY FILE|>")
	} else {
		lines := strings.Split(printableFileContents, "\n")
		for _, line := range lines {
			l.Plainf("%s", line)
		}
	}
}
