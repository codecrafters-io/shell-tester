package internal

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testA4(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	executableName := "custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999))
	logger.UpdateSecondaryPrefix("setup")
	logger.Infof("Available executables:\n- %s", executableName)
	logger.ResetSecondaryPrefix()
	_, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: executableName, CommandMetadata: getRandomString()},
	}, true)
	if err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := "custom"
	completion := executableName
	completionEndsWithNoSpace := false

	err = test_cases.CommandAutocompleteTestCase{
		RawCommand:         command,
		ExpectedReflection: completion,
		SuccessMessage:     fmt.Sprintf("Received completion for %q", command),
		ExpectedAutocompletedReflectionHasNoSpace: completionEndsWithNoSpace,
		SkipPromptAssertion:                       true,
	}.Run(asserter, shell, logger, false)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
