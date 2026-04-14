package internal

import (
	"fmt"
	"path"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type partialAndCompleteAutocompletePair struct {
	partial  string
	complete string
}

func testPA6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	completerDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	command := "git"
	subCommand := "remote"
	completerPath := path.Join(completerDir, "gitRemoteCompleter")

	choices := []partialAndCompleteAutocompletePair{
		{partial: "set", complete: "set-url"},
		{partial: "get", complete: "get-url"},
		{partial: "ad", complete: "add"},
		{partial: "ren", complete: "rename"},
	}

	choice := choices[random.RandomInt(0, len(choices))]

	if err := (&custom_executable.CompleterExecutableSpecification{
		Path:        completerPath,
		SecretValue: getRandomString(),
		CompleterConfiguration: completer_configuration.CompleterConfiguration{
			OutputLines: []string{choice.complete},
			ExpectedArguments: &completer_configuration.CompleterConfigurationExpectedArguments{
				Argv1: command,
				Argv2: choice.partial,
				Argv3: subCommand,
			},
		},
	}).Create(); err != nil {
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

	autocompleteTestCase := test_cases.AutocompleteTestCase{
		RawInput:            fmt.Sprintf("%s %s %s", command, subCommand, choice.partial),
		ExpectedCompletion:  fmt.Sprintf("%s %s %s ", command, subCommand, choice.complete),
		SkipPromptAssertion: true,
	}
	if err := autocompleteTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
