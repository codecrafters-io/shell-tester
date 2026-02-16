package internal

import (
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
	stageLogger := stageHarness.Logger

	inputArgsAndCompletion := []inputArgsAndCompletion{
		{Input: "ech", Completion: "echo", CompletionEndsWithNoSpace: false, Args: []string{"hello"}, Response: "hello"},
		{Input: "ech", Completion: "echo", CompletionEndsWithNoSpace: false, Args: []string{"foo", "bar"}, Response: "foo bar"},
	}

	for _, inputArgsAndCompletion := range inputArgsAndCompletion {
		err := a2Helper(stageHarness, stageLogger, inputArgsAndCompletion.Input, inputArgsAndCompletion.Completion, inputArgsAndCompletion.CompletionEndsWithNoSpace, inputArgsAndCompletion.Args, inputArgsAndCompletion.Response)
		if err != nil {
			return err
		}
		stageLogger.Infof("Tearing down shell")
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
		InputText:          command,
		ExpectedReflection: completion,
		ExpectedAutocompletedReflectionHasNoSpace: completionEndsWithNoSpace,
		Args:                args,
		ExpectedOutput:      response,
		SkipPromptAssertion: true,
	}.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
