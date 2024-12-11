package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/screen_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	if err := shell.Start(); err != nil {
		return err
	}

	screenAsserter := screen_asserter.NewScreenAsserter(shell, logger)
	if err := screenAsserter.Shell.ReadUntil(AsBool(screenAsserter.RunWithPromptAssertion)); err != nil {
		if err := screenAsserter.RunWithPromptAssertion(); err != nil {
			return err
		}
	}

	commandResponseTestCase := test_cases.NewCommandResponseTestCase("nonexistent", "bash: nonexistent: command not found", nil, "")
	if err := commandResponseTestCase.Run(screenAsserter); err != nil {
		return err
	}

	logger.Successf("$ ")
	return nil
}
