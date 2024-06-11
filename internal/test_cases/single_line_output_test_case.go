package test_cases

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
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
		shell.LogOutput(cleanOutput(output))
	}

	if err != nil {
		// Here, we are sure we have read the entire output, so we don't read any more
		if errors.Is(err, shell_executable.ErrConditionNotMet) {
			return fmt.Errorf("Expected first line of output to end with '\\n' (newline), got %q", string(cleanOutput(output)))
		} else if errors.Is(err, shell_executable.ErrProgramExited) {
			return fmt.Errorf("Expected shell to be a long-running process, but it exited")
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
		// If test fails, we still want to log the rest of the output
		restOfOutput, err := shell.ReadBytesUntilTimeout(100 * time.Millisecond)
		if err == nil {
			shell.LogOutput(cleanOutput(restOfOutput))
		}
		return fmt.Errorf("Expected first line of output to %s, got %q", t.ExpectedPatternExplanation, string(cleanOutput(output)))
	}

	logger.Successf("✓ %s", t.SuccessMessage)

	return nil
}

func stripSpaceRight(buf []byte) []byte {
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

func cleanOutput(buf []byte) []byte {
	buf = shell_executable.StripANSI(buf)
	buf = reversePTYTransformation(buf)
	buf = stripSpaceRight(buf)
	return buf
}
