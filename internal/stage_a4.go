package internal

import (
	"fmt"
	"os"
	"path"
	"strconv"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
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

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	// // // // //
	executableDir, err := getRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	// Add executableDir to PATH (That is where the my_exe file is created)
	currentPath := os.Getenv("PATH")
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", executableDir, currentPath))

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	randomCode := getRandomString()
	randomExecutableName := "custom_exe_" + strconv.Itoa(random.RandomInt(1000, 9999))
	exePath := path.Join(executableDir, randomExecutableName)

	err = custom_executable.CreateSignaturePrinterExecutable(randomCode, exePath)
	if err != nil {
		return err
	}
	// // // // //

	command := "custom"
	completion := randomExecutableName
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

	logger.Infof("Tearing down shell")
	return logAndQuit(asserter, nil)
}
