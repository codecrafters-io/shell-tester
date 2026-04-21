package internal

import (
	"fmt"
	"path"
	"regexp"
	"strconv"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA10(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	completerDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	commandName := "git"
	completerPath := path.Join(completerDir, "singleCompleter")

	completionSubcommand := random.RandomElementFromArray(
		[]string{"clone", "add", "commit", "push"},
	)

	compLine := commandName + " "
	compPoint := strconv.Itoa(len(compLine))

	// We need not use a valid completer here, but a valid completer makes it
	// obvious in the error messages
	// that the completion script was run when it shouldn't have
	if err := (&custom_executable.CompleterExecutableSpecification{
		Path:        completerPath,
		SecretValue: getRandomString(),
		CompleterConfiguration: completer_configuration.CompleterConfiguration{
			OutputLines: []string{completionSubcommand},
			ExpectedArguments: &completer_configuration.CompleterConfigurationExpectedArguments{
				Argv1: commandName,
				Argv2: "",
				Argv3: commandName,
			},
			ExpectedEnvVars: &completer_configuration.CompleterConfigurationExpectedEnvVars{
				CompLine:  compLine,
				CompPoint: compPoint,
			},
		},
	}).Create(); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	registerCmd := fmt.Sprintf("complete  -C  '%s'  %s", completerPath, commandName)

	if err := (test_cases.CommandWithNoResponseTestCase{
		Command:        registerCmd,
		SuccessMessage: "✓ Registered command-based completion",
	}).Run(asserter, shell, logger, false); err != nil {
		return err
	}

	unregisterCmd := fmt.Sprintf("complete -r %s", commandName)

	if err := (test_cases.CommandWithNoResponseTestCase{
		Command:        unregisterCmd,
		SuccessMessage: "✓ No output from complete -r",
	}).Run(asserter, shell, logger, false); err != nil {
		return err
	}

	printCompletionSpecTestCase := test_cases.CommandResponseTestCase{
		Command:        fmt.Sprintf("complete -p %s", commandName),
		ExpectedOutput: fmt.Sprintf("complete: %s: no completion specification", commandName),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^bash: complete: %s: no completion specification$`, regexp.QuoteMeta(commandName))),
		},
		SuccessMessage: "✓ Found missing completion specification after unregister",
	}
	if err := printCompletionSpecTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	if err := (test_cases.AutocompleteTestCase{
		RawInput:            commandName + " ",
		ExpectedCompletion:  commandName + " ",
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
