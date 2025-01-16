package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testA6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	prefix := "xyz_"
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	executableNames := []string{}
	logger.UpdateSecondaryPrefix("setup")
	logger.Infof("Available executables:")
	// TODO: Not supposed to be here but can't think of how to do this cleanly
	for _, word := range randomWords {
		executableName := prefix + word
		prefix = executableName + "_"
		executableNames = append(executableNames, executableName)
		logger.Infof("- %s", executableName)
	}
	logger.ResetSecondaryPrefix()

	for _, executableName := range executableNames {
		_, err := SetUpSignaturePrinter(stageHarness, shell, getRandomString(), executableName)
		if err != nil {
			return err
		}
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	panic("stop")

	// command := prefix
	// sort.Strings(executableNames)
	// completions := strings.Join(executableNames, "  ")
	// completionEndsWithNoSpace := true

	// err := test_cases.CommandMultipleCompletionsTestCase{
	// 	RawCommand:         command,
	// 	TabCount:           2,
	// 	ExpectedReflection: completions,
	// 	SuccessMessage:     fmt.Sprintf("Received completion for %q", command),
	// 	ExpectedAutocompletedReflectionHasNoSpace: completionEndsWithNoSpace,
	// 	CheckForBell:        true,
	// 	SkipPromptAssertion: true,
	// }.Run(asserter, shell, logger, false)
	// if err != nil {
	// 	return err
	// }

	return logAndQuit(asserter, nil)
}
