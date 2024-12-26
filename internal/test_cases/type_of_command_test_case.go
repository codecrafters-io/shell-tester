package test_cases

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type TypeOfCommandTestCase struct {
	Command        string
	SuccessMessage string
}

func (t *TypeOfCommandTestCase) RunForBuiltin(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if t.SuccessMessage == "" {
		t.SuccessMessage = "✓ Received expected response"
	}

	testCase := CommandResponseTestCase{
		Command:          fmt.Sprintf("type %s", t.Command),
		ExpectedOutput:   fmt.Sprintf(`%s is a shell builtin`, t.Command),
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^%s is a( special)? shell builtin$`, t.Command))},
		SuccessMessage:   t.SuccessMessage,
	}

	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}

func (t *TypeOfCommandTestCase) RunForExecutable(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger, customExecutablePath string) error {
	var expectedPath string

	if t.Command == "my_exe" {
		expectedPath = customExecutablePath
	} else {
		path, err := exec.LookPath(t.Command)
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Error finding %s in PATH", t.Command)
		}
		expectedPath = path
	}

	testCase := CommandResponseTestCase{
		Command:          fmt.Sprintf("type %s", t.Command),
		ExpectedOutput:   fmt.Sprintf(`%s is %s`, t.Command, expectedPath),
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^(%s is )?%s$`, t.Command, expectedPath))},
		SuccessMessage:   "✓ Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}

func (t *TypeOfCommandTestCase) RunForInvalidCommand(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := CommandResponseTestCase{
		Command:          fmt.Sprintf("type %s", t.Command),
		ExpectedOutput:   fmt.Sprintf(`%s: not found`, t.Command),
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^(bash: type: )?%s[:]? not found$`, t.Command))},
		SuccessMessage:   "✓ Received expected response",
	}

	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return nil
}
