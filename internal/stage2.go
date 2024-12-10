package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	screenAsserter := assertions.NewScreenAsserter(shell, logger)

	// Checks if prompt is present
	if err := screenAsserter.Run(); err != nil {
		return err
	}

	// TODO: Can shorten into a SingleLineCommandTestCase
	// ------ Test case starts
	shell.SendCommand("nonexistent")
	screenAsserter.PushAssertion(screenAsserter.SingleLineAssertion(0, "$ nonexistent", nil, "nonexistent"))
	screenAsserter.PushAssertion(screenAsserter.SingleLineAssertion(1, "", nil, "nonexistent: command not found"))

	if err := screenAsserter.Run(); err != nil {
		return err
	}
	// ------ Test case ends

	// if err := screenAsserter.RunWithoutLastPromptAssertion(); err != nil {
	// 	return err
	// }

	// [x] Assert prompt is printed
	// [x] Send command_1
	// [x] Assert "$ command_1" is present
	// [] Assert next line "command_1: not found" is present
	// [] Assert prompt is printed again
	// [] Send command_2
	// [] Assert "$ command_2" is present
	// [] Assert next line "command_2: not found" is present
	// [] Assert prompt is printed again

	// screenAsserter := assertions.NewScreenAsserter(shell, logger)
	// promptAssertion := screenAsserter.PromptAssertion(0, "$ ")
	// screenAsserter.PushAssertion(&promptAssertion)

	// responseTestCase := test_cases.NewResponseTestCase()

	// if err := responseTestCase.Run(screenAsserter, true); err != nil {
	// 	return err
	// }

	// screenAsserter.ClearAssertions()
	// firstLineAssertion := screenAsserter.SingleLineAssertion(0, "$ nonexistent", nil, "nonexistent")
	// screenAsserter.PushAssertion(&firstLineAssertion)
	// commandResponseTestCase := test_cases.NewCommandResponseTestCase("nonexistent")
	// if err := commandResponseTestCase.Run(screenAsserter, true); err != nil {
	// 	return err
	// }

	// secondLineAssertion := screenAsserter.SingleLineAssertion(1, "", []*regexp.Regexp{regexp.MustCompile(`bash: nonexistent: command not found`)}, "nonexistent: command not found")
	// screenAsserter.PushAssertion(&secondLineAssertion)

	// // At this stage the user might or might not have implemented a REPL to print the prompt again, so we won't test further
	// // ToDo: Remove this prompt assertion from here
	// promptAssertion = screenAsserter.PromptAssertion(2, "$ ")
	// screenAsserter.PushAssertion(&promptAssertion)

	// if err := responseTestCase.Run(screenAsserter, true); err != nil {
	// 	return err
	// }

	return nil
}
