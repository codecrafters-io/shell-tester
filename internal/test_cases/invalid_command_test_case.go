package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type InvalidCommandTestCase struct {
	Command string
}

func (t *InvalidCommandTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := CommandResponseTestCase{
		Command:          t.Command,
		ExpectedOutput:   t.getExpectedOutput(),
		FallbackPatterns: t.getFallbackPatterns(),
		SuccessMessage:   "✓ Received command not found message",
	}

	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}

func (t *InvalidCommandTestCase) RunWithoutNextPromptAssertion(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := CommandReflectionTestCase{
		Command:             t.Command,
		SkipPromptAssertion: true,
	}
	if err := testCase.Run(asserter, shell, logger, true); err != nil {
		return err
	}

	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput:   t.getExpectedOutput(),
		FallbackPatterns: t.getFallbackPatterns(),
	})

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	return nil
}

func (t *InvalidCommandTestCase) getExpectedOutput() string {
	return fmt.Sprintf("%s: command not found", t.Command)
}

func (t *InvalidCommandTestCase) getFallbackPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`^(bash: )?` + t.Command + `: (command )?not found$`),
		regexp.MustCompile(`^ash: ` + t.Command + `: not found$`),
		regexp.MustCompile(`^zsh: command not found: ` + t.Command + `$`),
	}
}
