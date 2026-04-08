package internal

import (
	"fmt"
	"path"
	"strconv"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA9(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	command := "git"
	completerPath := path.Join(randomDir, "lcpEnvCompleter")
	secret := getRandomString()

	ambiguousMatches := []string{"cherry-pick", "checkout"}

	prepareCompleterExecutable := func(argv1, argv2, argv3, compLine string, outputLines []string) error {
		return (&custom_executable.CompleterExecutableSpecification{
			Path:        completerPath,
			SecretValue: secret,
			CompleterConfiguration: completer_configuration.CompleterConfiguration{
				OutputLines: outputLines,
				ExpectedArguments: &completer_configuration.CompleterConfigurationExpectedArguments{
					Argv1: argv1,
					Argv2: argv2,
					Argv3: argv3,
				},
				ExpectedEnvVars: &completer_configuration.CompleterConfigurationEnvVars{
					CompLine:  compLine,
					CompPoint: strconv.Itoa(len(compLine)),
				},
			},
		}).Create()
	}

	// Same shape as PA6/PA7: typed partial → completion segment on the line after TAB.
	steps := []partialAndCompleteAutocompletePair{
		{partial: "c", complete: "che"},
		{partial: "r", complete: "cherry-pick "},
	}

	// Line "git c": prefix "c" → LCP of checkout / cherry-pick is "che".
	if err := prepareCompleterExecutable(
		command,
		steps[0].partial,
		command,
		fmt.Sprintf("%s %s", command, steps[0].partial),
		ambiguousMatches,
	); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	registerCmd := fmt.Sprintf("complete -C %s %s", completerPath, command)
	registerTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        registerCmd,
		SuccessMessage: "✓ Registered command-based completion",
	}
	if err := registerTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	pairs := []test_cases.InputAndCompletionPair{
		{
			Input:              fmt.Sprintf("%s %s", command, steps[0].partial),
			ExpectedCompletion: fmt.Sprintf("%s %s", command, steps[0].complete),
			CallbackAfterAutocompletion: func() error {
				// Line will be "git cher" before the final TAB — unique match cherry-pick (vs checkout's "chec…").
				return prepareCompleterExecutable(
					command,
					fmt.Sprintf("%s%s", steps[0].complete, steps[1].partial),
					command,
					fmt.Sprintf("%s %s%s", command, steps[0].complete, steps[1].partial),
					[]string{"cherry-pick"},
				)
			},
		},
		{
			Input:              steps[1].partial,
			ExpectedCompletion: fmt.Sprintf("%s %s", command, steps[1].complete),
		},
	}

	if err := (test_cases.PartialCompletionsTestCase{
		InputAndCompletionPairs: pairs,
		SuccessMessage:          "✓ Received all partial completions",
		SkipPromptAssertion:     true,
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
