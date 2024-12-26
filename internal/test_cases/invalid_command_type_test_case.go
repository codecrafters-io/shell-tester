package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type InvalidCommandTypeTestCase struct {
	Command string
}

func (t *InvalidCommandTypeTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := CommandResponseTestCase{
		Command:          "type " + t.Command,
		ExpectedOutput:   t.getExpectedOutput(),
		FallbackPatterns: t.getFallbackPatterns(),
		SuccessMessage:   "✓ Received expected response",
	}

	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}

func (t *InvalidCommandTypeTestCase) getExpectedOutput() string {
	return fmt.Sprintf("%s: not found", t.Command)
}

func (t *InvalidCommandTypeTestCase) getFallbackPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^(bash: type: )?%s[:]? not found$`, t.Command))}
}
