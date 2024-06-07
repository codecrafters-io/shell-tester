package test_cases

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// SingleLineOutputTestCase verifies a prompt exists, sends a command and matches the output against a string.
type SingleLineOutputTestCase struct {
	// The command to execute (the command's output will be matched against ExpectedPattern)
	Command string

	// TODO: Do we still need a regex here?
	// ExpectedPattern is the regex that is evaluated against the command's output.
	// Add \r\n at the end of the pattern if you're expecting a newline.
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
	// fmt.Printf("output: %q\n", output)

	// Whether the condition fails on not, we want to log the output
	if len(output) > 0 {
		trimmedOutput := bytes.TrimRightFunc(output, func(r rune) bool {
			return r == '\r' || r == '\n'
		})
		shell.LogOutput(shell_executable.StripANSI(trimmedOutput))
	}

	if err != nil {
		if errors.Is(err, shell_executable.ErrConditionNotMet) {
			logger.Errorf("Expected first line of output to end with '\\n' (newline), got %q", string(shell_executable.StripANSI(output)))
		} else if errors.Is(err, shell_executable.ErrProgramExited) {
			logger.Errorf("Expected shell to be a long-running process, but it exited")
		}
		return err
	}

	// We will explicitly support the case where user uses CR/LF at the end
	// and we receive CR/CR/LF
	if len(output) > 2 && bytes.Equal(output[len(output)-3:], []byte{'\r', '\r', '\n'}) {
		output = append(output[:len(output)-3], []byte{'\r', '\n'}...)
	}
	output = shell_executable.StripANSI(output)

	if !t.ExpectedPattern.Match(output) {
		logger.Errorf("Expected first line of output to %s, got %q", t.ExpectedPatternExplanation, string(TrimRightSpace(output)))
		return fmt.Errorf("Output did not match expected pattern")
	}

	logger.Successf("✓ %s", t.SuccessMessage)

	return nil
}

func TrimRightSpace(buf []byte) []byte {
	return bytes.TrimRightFunc(buf, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
}
