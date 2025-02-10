package internal

import (
	"os"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	executableName := "my_exe"

	// We create 2 custom executables with the same name
	// in 2 different directories
	// The one in the first directory has its executable bit removed
	// The other one is in a directory which comes after the first one's directory on PATH
	// The other one should still be found & used
	executable_expected_to_be_not_found := "my_exe"
	executableDir, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: executable_expected_to_be_not_found, CommandMetadata: getRandomString()},
	}, true)
	if err != nil {
		return err
	}
	notExePath := filepath.Join(executableDir, executable_expected_to_be_not_found)
	currentPerms, _ := os.Stat(notExePath)
	os.Chmod(notExePath, currentPerms.Mode() & ^os.FileMode(0o111))

	executableDir, err = SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: executableName, CommandMetadata: getRandomString()},
	}, true)
	if err != nil {
		return err
	}

	logAvailableExecutables(logger, []string{executableName})
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

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
			expectedPath = filepath.Join(executableDir, executableName)
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
