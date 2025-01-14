package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type inputArgsAndCompletion struct {
	Input                     string
	Completion                string
	CompletionEndsWithNoSpace bool
	Args                      []string
	Response                  string
}

func testA2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger

	inputArgsAndCompletion := []inputArgsAndCompletion{
		{Input: "ech", Completion: "echo", CompletionEndsWithNoSpace: false, Args: []string{"hello"}, Response: "hello"},
		{Input: "typ", Completion: "type", CompletionEndsWithNoSpace: true, Args: []string{" type"}, Response: "type is a shell builtin"},
	}

	for _, inputArgsAndCompletion := range inputArgsAndCompletion {
		err := a2Helper(stageHarness, logger, inputArgsAndCompletion.Input, inputArgsAndCompletion.Completion, inputArgsAndCompletion.CompletionEndsWithNoSpace, inputArgsAndCompletion.Args, inputArgsAndCompletion.Response)
		if err != nil {
			return err
		}
	}

	return nil
}

func a2Helper(stageHarness *test_case_harness.TestCaseHarness, logger *logger.Logger, command string, completion string, completionEndsWithNoSpace bool, args []string, response string) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	err := test_cases.CommandAutocompleteAndResponseTestCase{
		RawCommand:         command,
		ExpectedReflection: completion,
		ExpectedAutocompletedReflectionHasNoSpace: completionEndsWithNoSpace,
		Args:                args,
		ExpectedOutput:      response,
		SuccessMessage:      fmt.Sprintf("Received completion for %q", command),
		SkipPromptAssertion: true,
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	logger.Infof("Tearing down shell")
	return logAndQuit(asserter, nil)
}
