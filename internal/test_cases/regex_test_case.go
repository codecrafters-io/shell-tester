package test_cases

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// RegexTestCase verifies a prompt exists, sends a command and matches the output against a regex pattern.
type RegexTestCase struct {
	// The command to execute (the command's output will be matched against ExpectedPattern)
	Command string

	// ExpectedPattern is the regex that is evaluated against the command's output.
	// Add \r\n at the end of the pattern if you're expecting a newline.
	ExpectedPattern *regexp.Regexp

	// ExpectedPatternExplanation is used in the error message if the ExpectedPattern doesn't match the command's output
	ExpectedPatternExplanation string

	// SuccessMessage is logged if the ExpectedPattern matches the command's output
	SuccessMessage string
}

func (t RegexTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	promptTestCase := NewSilentPromptTestCase("$ ")

	if err := promptTestCase.Run(shell, logger); err != nil {
		return err
	}

	if err := shell.SendCommand(t.Command); err != nil {
		return err
	}

	regexMatchCondition := func(buf []byte) bool {
		return t.ExpectedPattern.Match(shell_executable.StripANSI(buf))
	}

	output, err := shell.ReadBytesUntil(regexMatchCondition)

	// Whether the condition fails on not, we want to log the output
	if len(output) > 0 {
		// TODO: Avoid this clunkiness + avoid "\r\n" in error message
		if string(output[len(output)-2:]) == "\r\n" {
			shell.LogOutput(shell_executable.StripANSI(output[:len(output)-2]))
		} else {
			shell.LogOutput(shell_executable.StripANSI(output))
		}
	}

	if err != nil {
		if err == shell_executable.ErrConditionNotMet {
			logger.Errorf("Expected output to %s, got %q", t.ExpectedPatternExplanation, string(shell_executable.StripANSI(output)))
		}

		return err
	}

	logger.Successf("✓ %s", t.SuccessMessage)

	return nil
}
