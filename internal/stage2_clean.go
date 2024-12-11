package internal

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	if err := shell.Start(); err != nil {
		return err
	}

	screenAsserter := assertions.NewScreenAsserter(shell, logger)
	if err := screenAsserter.Shell.ReadUntil(AsBool(screenAsserter.RunWithPromptAssertion)); err != nil {
		return err
	}

	shell.SendCommand("nonexistent")
	screenAsserter.PushAssertion(screenAsserter.SingleLineAssertion(0, "$ nonexistent", nil, "nonexistent"))
	screenAsserter.PushAssertion(screenAsserter.SingleLineAssertion(1, "", []*regexp.Regexp{regexp.MustCompile(`^bash: nonexistent: command not found$`)}, "bash: nonexistent: command not found"))
	if err := screenAsserter.Shell.ReadUntil(AsBool(screenAsserter.RunWithPromptAssertion)); err != nil {
		return err
	}

	logger.Successf("$ ")
	return nil
}
