package test_cases

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// OutputOnlyTestCase should be used when no input is to be sent to the shell
// but output is expected (due to a side effect of writing to a file) from a running foreground/background process
type OutputOnlyTestCase struct {
	ExpectedOutputLines []string
	SuccessMessage      string
}

func (t OutputOnlyTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if len(t.ExpectedOutputLines) == 0 {
		panic("Codecrafters Internal Error - ExpectedOutputLines is empty in BackgroundCommandOutputOnlyTestCase")
	}

	outputLinesAssertion := assertions.NewMultiLineAssertion(t.ExpectedOutputLines)
	asserter.AddAssertion(&outputLinesAssertion)

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	if t.SuccessMessage != "" {
		logger.Successf("%s", t.SuccessMessage)
	}

	return nil
}
