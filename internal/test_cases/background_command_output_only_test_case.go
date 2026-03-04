package test_cases

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type BackgroundCommandOutputOnlyTestCase struct {
	ExpectedOutputLines []string
	SuccessMessage      string
}

func (t BackgroundCommandOutputOnlyTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	if len(t.ExpectedOutputLines) == 0 {
		panic("Codecrafters Internal Error - ExpectedOutputLines is empty in BackgroundCommandOutputOnlyTestCase")
	}

	outputInPromptLine := t.ExpectedOutputLines[0]

	promptLineReflection := fmt.Sprintf("$ %s", outputInPromptLine)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: promptLineReflection,
	})

	remainingOutputLinesAssertion := assertions.NewMultiLineAssertion(t.ExpectedOutputLines[1:])
	asserter.AddAssertion(&remainingOutputLinesAssertion)

	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	if t.SuccessMessage != "" {
		logger.Successf("%s", t.SuccessMessage)
	}

	return nil
}
