package internal

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	screenAsserter := assertions.NewScreenAsserter(shell, logger)
	promptAssertion := screenAsserter.PromptAssertion(0, "$ ", screenAsserter)
	screenAsserter.AddAssertion(&promptAssertion)

	responseTestCase := test_cases.NewResponseTestCase()

	if err := responseTestCase.Run(screenAsserter, true); err != nil {
		return err
	}

	screenAsserter.ClearAssertions()
	firstLineAssertion := screenAsserter.SingleLineAssertion(0, "$ nonexistent", nil, "nonexistent", screenAsserter)
	screenAsserter.AddAssertion(&firstLineAssertion)
	commandResponseTestCase := test_cases.NewCommandResponseTestCase("nonexistent")
	if err := commandResponseTestCase.Run(screenAsserter, true); err != nil {
		return err
	}

	secondLineAssertion := screenAsserter.SingleLineAssertion(1, "", []*regexp.Regexp{regexp.MustCompile(`bash: nonexistent: command not found`)}, "nonexistent: command not found", screenAsserter)
	screenAsserter.AddAssertion(&secondLineAssertion)

	// At this stage the user might or might not have implemented a REPL to print the prompt again, so we won't test further
	// ToDo: Remove this prompt assertion from here
	promptAssertion = screenAsserter.PromptAssertion(2, "$ ", screenAsserter)
	screenAsserter.AddAssertion(&promptAssertion)

	if err := responseTestCase.Run(screenAsserter, true); err != nil {
		return err
	}

	return nil
}
