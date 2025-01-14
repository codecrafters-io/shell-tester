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
	logger := stageHarness.Logger

	inputAndCompletion := []inputAndCompletion{
		{Input: "typ", Completion: "type", CompletionEndsWithNoSpace: true},
		{Input: "ech", Completion: "echo", CompletionEndsWithNoSpace: false},
		{Input: "exi", Completion: "exit", CompletionEndsWithNoSpace: false},
	}

	for _, inputAndCompletion := range inputAndCompletion {
		err := a1Helper(stageHarness, logger, inputAndCompletion.Input, inputAndCompletion.Completion, inputAndCompletion.CompletionEndsWithNoSpace)
		if err != nil {
			return err
		}
	}

	return nil
}

func a1Helper(stageHarness *test_case_harness.TestCaseHarness, logger *logger.Logger, command string, completion string, completionEndsWithNoSpace bool) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(); err != nil {
		return err
	}

	err := test_cases.CommandAutocompleteTestCase{
		RawCommand:         command,
		ExpectedReflection: completion,
		SuccessMessage:     fmt.Sprintf("Received completion for %q", command),
		ExpectedAutocompletedReflectionHasNoSpace: completionEndsWithNoSpace,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	logger.Infof("Tearing down shell")
	return logAndQuit(asserter, nil)
}
