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
	"github.com/codecrafters-io/tester-utils/random"
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

	prepareCompleterExecutable := func(argv1, argv2, argv3, compLineEnvVar string, outputLines []string) error {
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
					CompLine:  compLineEnvVar,
					CompPoint: strconv.Itoa(len(compLineEnvVar)),
				},
			},
		}).Create()
	}

	lcpBranches := []struct {
		disambiguationLetter string
		subcommandName       string
	}{
		{disambiguationLetter: "c", subcommandName: "checkout"},
		{disambiguationLetter: "r", subcommandName: "cherry-pick"},
	}
	chosenBranch := lcpBranches[random.RandomInt(0, len(lcpBranches))]

	sortedAmbiguousSubcommands := []string{"checkout", "cherry-pick"}

	partialAndCompletePairs := []partialAndCompleteAutocompletePair{
		{partial: "c", complete: "che"},
		{partial: chosenBranch.disambiguationLetter, complete: chosenBranch.subcommandName + " "},
	}

	// git c -> git che (LCP Completion)
	if err := prepareCompleterExecutable(
		command,
		partialAndCompletePairs[0].partial,
		command,
		fmt.Sprintf("%s %s", command, partialAndCompletePairs[0].partial),
		sortedAmbiguousSubcommands,
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

	partialCompletionPairs := []test_cases.InputAndCompletionPair{
		{
			Input:              fmt.Sprintf("%s %s", command, partialAndCompletePairs[0].partial),
			ExpectedCompletion: fmt.Sprintf("%s %s", command, partialAndCompletePairs[0].complete),
			// After the first tab completion, replace the completer script with a new script that expects
			// the updated argv and env variables
			CallbackAfterAutocompletion: func() error {
				secondTabPartialWord := fmt.Sprintf("%s%s", partialAndCompletePairs[0].complete, partialAndCompletePairs[1].partial)
				compLineEnvVarBeforeSecondTab := fmt.Sprintf("%s %s%s", command, partialAndCompletePairs[0].complete, partialAndCompletePairs[1].partial)
				return prepareCompleterExecutable(
					command,
					secondTabPartialWord,
					command,
					compLineEnvVarBeforeSecondTab,
					[]string{chosenBranch.subcommandName},
				)
			},
		},
		{
			Input:              partialAndCompletePairs[1].partial,
			ExpectedCompletion: fmt.Sprintf("%s %s", command, partialAndCompletePairs[1].complete),
		},
	}

	if err := (test_cases.PartialCompletionsTestCase{
		InputAndCompletionPairs: partialCompletionPairs,
		SuccessMessage:          "✓ Received all partial completions",
		SkipPromptAssertion:     true,
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
