package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	myExeCommandName := "my_exe"

	// Test PATH resolution with duplicate executable names
	//
	// This test creates three files with identical names ("my_exe") in different directories:
	// - myExe3 with execute permissions removed
	// - myExe2 with normal permissions
	// - myExe1 with execute permissions removed
	// Since we prepend to PATH, it will look like myExe1:myExe2:myExe3:...
	//
	// Expected behavior:
	// - When the command is executed, the shell should skip myExe1 (not executable)
	// - The shell should continue searching PATH and find/execute myExe2
	// - The purpose of myExe3 is to catch a wrong solution which traverses PATH in reverse
	// - This verifies proper PATH traversal and permission checking

	// myExe3
	if _, err := setUpNonExecutable(stageHarness, shell, myExeCommandName); err != nil {
		return err
	}

	// myExe2
	executableDir, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: myExeCommandName, CommandMetadata: getRandomString()},
	}, true)
	if err != nil {
		return err
	}

	// myExe1
	nonExePath, err := setUpNonExecutable(stageHarness, shell, myExeCommandName)
	if err != nil {
		return err
	}

	logPath(shell, logger, 36) // Prefix length is 36 characters for this stage
	logAvailableExecutables(logger, []string{myExeCommandName})

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
			expectedPath = filepath.Join(executableDir, myExeCommandName)

			// Alpine Busybox has a bug where it doesn't check permissions
			if isTestingTesterUsingBusyboxOnAlpine(stageHarness) {
				expectedPath = nonExePath
			}
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

func logPath(shell *shell_executable.ShellExecutable, logger *logger.Logger, prefixLength int) {
	path := shell.GetPath()
	lengthLimit := 80 - prefixLength

	if len(path) > lengthLimit {
		pathChunks := strings.Split(path, ":")
		path = ""

		for _, chunk := range pathChunks {
			// 4 is reserved for the colon and the ellipsis
			if len(path)+len(chunk) <= lengthLimit-4 {
				path += chunk + ":"
			} else {
				path += "..."
				break
			}
		}
	}

	logger.UpdateLastSecondaryPrefix("setup")
	logger.Infof("PATH is now: %s", path)
	logger.ResetSecondaryPrefixes()
}

func setUpNonExecutable(stageHarness *test_case_harness.TestCaseHarness, shell *shell_executable.ShellExecutable, commandName string) (string, error) {
	nonExeDir, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: commandName, CommandMetadata: getRandomString()},
	}, true)
	if err != nil {
		return "", err
	}

	nonExePath := filepath.Join(nonExeDir, commandName)
	currentPerms, _ := os.Stat(nonExePath)
	os.Chmod(nonExePath, currentPerms.Mode() & ^os.FileMode(0o111))

	return nonExePath, nil
}
