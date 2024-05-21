package test_cases

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type RegexTestCase struct {
	Command                    string
	ExpectedPattern            *regexp.Regexp
	ExpectedPatternExplanation string
	SuccessMessage             string
}

func (t RegexTestCase) Run(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if err := shell.AssertPrompt("$ "); err != nil {
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
		shell.LogOutput(output)
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
