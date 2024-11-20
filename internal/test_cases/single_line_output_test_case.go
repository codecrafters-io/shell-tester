package test_cases

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/fatih/color"
)

// SingleLineOutputTestCase verifies a prompt exists, sends a command and matches the output against a string.
type SingleLineOutputTestCase struct {
	// The command to execute (the command's output will be matched against ExpectedPattern)
	Command string

	// ExpectedPattern is the regex that is evaluated against the command's output.
	ExpectedPattern *regexp.Regexp

	// ExpectedPatternExplanation is used in the error message if the ExpectedPattern doesn't match the command's output
	ExpectedPatternExplanation string

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t SingleLineOutputTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	promptTestCase := NewSilentPromptTestCase("$ ")

	if err := promptTestCase.Run(shell, logger); err != nil {
		return err
	}

	if err := shell.SendCommand(t.Command); err != nil {
		return err
	}

	CRLFCondition := func(buf []byte) bool {
		return len(buf) > 1 && bytes.Equal(buf[len(buf)-2:], []byte{'\r', '\n'})
	}

	output, err := shell.ReadBytesUntil(CRLFCondition)

	// Whether the condition fails on not, we want to log the output
	if len(output) > 0 {
		shell.LogOutput(sanitizeLogOutput(output))
	}

	cleanedOutput := sanitizeLogOutput(output)
	if err != nil {
		// Here, we are sure we have read the entire output, so we don't read any more
		if errors.Is(err, shell_executable.ErrConditionNotMet) {
			return fmt.Errorf("Expected first line of output to end with '\\n' (newline), got %q", string(cleanedOutput))
		} else if errors.Is(err, shell_executable.ErrProgramExited) {
			exitCode := shell.ExitCode()
			if exitCode == -1 {
				return fmt.Errorf("Expected first line of output to end with '\\n' (newline), got %q", string(cleanedOutput))
			} else {
				return fmt.Errorf("Expected first line of output to end with '\\n' (newline), got %q. Program exited with code %d", string(cleanedOutput), exitCode)
			}
		}
	}

	if !t.ExpectedPattern.Match(cleanedOutput) {
		// If test fails, we still want to log the rest of the output
		restOfOutput, err := shell.ReadBytesUntilTimeout(100 * time.Millisecond)
		if err == nil {
			shell.LogOutput(sanitizeLogOutput(restOfOutput))
		}
		return fmt.Errorf(getColoredErrorMessage(t.ExpectedPatternExplanation, string(cleanedOutput)))
	}

	logger.Successf("✓ %s", t.SuccessMessage)

	return nil
}

func stripLineEnding(buf []byte) []byte {
	return bytes.TrimRightFunc(buf, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
}

func reversePTYTransformation(buf []byte) []byte {
	// PTY converts all bare LFs to CRLFs, so we just reverse that here
	// Bare CRs are left as is
	buf = bytes.ReplaceAll(buf, []byte{'\r', '\n'}, []byte{'\n'})
	return buf
}

func squashMultipleCR(buf []byte) []byte {
	// Squash multiple CRs into one
	re := regexp.MustCompile(`\r+`)
	return re.ReplaceAll(buf, []byte{'\r'})
}

func sanitizeLogOutput(buf []byte) []byte {
	buf = shell_executable.StripANSI(buf)
	buf = reversePTYTransformation(buf)
	buf = squashMultipleCR(buf)
	buf = stripLineEnding(buf)
	return buf
}

func colorizeString(colorToUse color.Attribute, msg string) string {
	c := color.New(colorToUse)
	return c.Sprint(msg)
}

func getColoredErrorMessage(expectedPatternExplanation string, cleanedOutput string) string {
	indent := 32 - 3 + len(strings.Split(expectedPatternExplanation, " ")[0]) + 1

	errorMsg := "Expected first line of output to" // 32
	errorMsg += " " + colorizeString(color.FgGreen, expectedPatternExplanation)
	errorMsg += "\n"
	errorMsg += strings.Repeat(" ", indent)
	errorMsg += colorizeString(color.FgRed, "got") // 3
	errorMsg += " " + colorizeString(color.FgRed, cleanedOutput)

	return errorMsg
}
