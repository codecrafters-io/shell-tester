package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/fatih/color"
)

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

func AsBool(T func() error) func() bool {
	// Takes in a function that takes no params & returns an error
	// Returns the function wrapped in a helper such that it returns a bool
	// in lieu of the error, true if the function execution is a success
	return func() bool { return T() == nil }
}

// LogReadableFileContents prints file contents in a readable way, replacing tabs and spaces with visible markers and coloring comments yellow.
func LogReadableFileContents(l *logger.Logger, fileContents string, logMsg string, fileName string) {
	l.Infof(logMsg)

	printableFileContents := strings.ReplaceAll(fileContents, "%", "%%")
	printableFileContents = strings.ReplaceAll(printableFileContents, "\t", "<|TAB|>")

	regex1 := regexp.MustCompile("[ ]+\n")
	regex2 := regexp.MustCompile("[ ]+$")
	printableFileContents = regex1.ReplaceAllString(printableFileContents, "\n")
	printableFileContents = regex2.ReplaceAllString(printableFileContents, "<|SPACE|>")

	if len(printableFileContents) == 0 {
		l.Plainf("[%s] <|EMPTY FILE|>", fileName)
	} else {
		lines := strings.Split(printableFileContents, "\n")
		// If the last line is empty (trailing newline), skip it
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		for _, line := range lines {
			if strings.Contains(line, "//") {
				code := strings.Split(line, "//")[0]
				comment := "//" + strings.Split(line, "//")[1]
				formattedLine := fmt.Sprintf("[%s] %s%s", fileName, code, ColorizeString(color.FgYellow, comment))
				l.Plainf(formattedLine)
			} else {
				l.Plainf("[%s] %s", fileName, line)
			}
		}
	}
}
