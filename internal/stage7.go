package internal

import (
	"fmt"
	"os"
	"path/filepath"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	// Add the random directory to PATH (where the my_exe file is created)
	randomDir, err := getRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	path := os.Getenv("PATH")
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, path))
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	customExecutablePath := filepath.Join(randomDir, "my_exe")
	err = custom_executable.CreateSignaturePrinterExecutable(getRandomString(), customExecutablePath)
	if err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	availableExecutables := []string{"cat", "cp", "mkdir", "my_exe"}

	for _, executable := range availableExecutables {
		testCase := test_cases.TypeOfCommandTestCase{
			Command: executable,
		}

		var expectedPath = ""
		if executable == "my_exe" {
			expectedPath = customExecutablePath
		}

		if err := testCase.RunForExecutable(asserter, shell, logger, expectedPath); err != nil {
			return err
		}
	}

	invalidCommands := getRandomInvalidCommands(2)

	for _, invalidCommand := range invalidCommands {
		testCase := test_cases.TypeOfCommandTestCase{
			Command: invalidCommand,
		}
		if err := testCase.RunForInvalidCommand(asserter, shell, logger); err != nil {
			return err
		}
	}

	return logAndQuit(asserter, nil)
}
