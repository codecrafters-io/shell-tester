package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type inputAndCompletion struct {
	Input                     string
	Completion                string
	CompletionEndsWithNoSpace bool
}

func testA1(stageHarness *test_case_harness.TestCaseHarness) error {
	stageLogger := stageHarness.Logger

	inputAndCompletions := []inputAndCompletion{
		{Input: "ech", Completion: "echo", CompletionEndsWithNoSpace: false},
		{Input: "exi", Completion: "exit", CompletionEndsWithNoSpace: false},
	}

	for _, inputAndCompletion := range inputAndCompletions {
		err := a1Helper(stageHarness, stageLogger, inputAndCompletion.Input, inputAndCompletion.Completion, inputAndCompletion.CompletionEndsWithNoSpace)
		if err != nil {
			return err
		}
		stageLogger.Infof("Tearing down shell")
	}

	return nil
}

func a1Helper(stageHarness *test_case_harness.TestCaseHarness, logger *logger.Logger, command string, completion string, completionEndsWithNoSpace bool) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	err := test_cases.CommandAutocompleteTestCase{
		RawCommand:         command,
		ExpectedReflection: completion,
		SuccessMessage:     fmt.Sprintf("Received completion for %q", command),
		ExpectedAutocompletedReflectionHasNoSpace: completionEndsWithNoSpace,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
