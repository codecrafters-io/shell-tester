package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testA6(stageHarness *test_case_harness.TestCaseHarness) error {
	stageLogger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	prefix := "xyz_"
	initialPrefix := prefix
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	executableNames := []string{}
	for _, word := range randomWords {
		executableName := prefix + word
		prefix = executableName + "_"
		_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
			{CommandType: "signature_printer", CommandName: executableName, CommandMetadata: getRandomString()},
		}, true)
		if err != nil {
			return err
		}
		executableNames = append(executableNames, executableName)
	}
	logAvailableExecutables(stageLogger, executableNames)

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	err := test_cases.PartialCompletionsTestCase{
		Inputs:              []string{initialPrefix, "_", "_"},
		ExpectedReflections: executableNames,
		SuccessMessage:      fmt.Sprintf("Received all partial completions for %q", executableNames[len(executableNames)-1]),
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageLogger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}

// TODO: Think of how to encapsulate this inside SetupExecutable function
func logAvailableExecutables(logger *logger.Logger, executableNames []string) {
	logger.UpdateLastSecondaryPrefix("setup")
	logger.Infof("Available executables:")
	for _, executableName := range executableNames {
		logger.Infof("- %s", executableName)
	}
	logger.ResetSecondaryPrefixes()
}
