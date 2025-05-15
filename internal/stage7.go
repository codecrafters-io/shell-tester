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
	executableName := "my_exe"

	// Test PATH resolution with duplicate executable names
	// This test creates two executables with identical names in different directories:
	// 1. e1: First executable in PATH, with execute permissions removed
	// 2. e2: Second executable in PATH, with normal permissions
	// Expected behavior:
	// - When the command is executed, the shell should skip e1 (not executable)
	// - The shell should continue searching PATH and find/execute e2
	// - This verifies proper PATH traversal and permission checking
	executableExpectedToNotBeFound := "my_exe"
	executableExpectedToNotBeFoundDir, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: executableExpectedToNotBeFound, CommandMetadata: getRandomString()},
	}, true)
	if err != nil {
		return err
	}
	notExePath := filepath.Join(executableExpectedToNotBeFoundDir, executableExpectedToNotBeFound)
	currentPerms, _ := os.Stat(notExePath)
	os.Chmod(notExePath, currentPerms.Mode() & ^os.FileMode(0o111))

	executableDir, err := SetUpCustomCommands(stageHarness, shell, []CommandDetails{
		{CommandType: "signature_printer", CommandName: executableName, CommandMetadata: getRandomString()},
	}, true)
	if err != nil {
		return err
	}

	logPath(shell, logger, 36) // Prefix length is 36 characters for this stage
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

func logPath(shell *shell_executable.ShellExecutable, logger *logger.Logger, prefixLength int) {
	path := shell.GetPath()
	lineLimit := 80 - prefixLength

	if len(path) > lineLimit {
		pathChunks := strings.Split(path, ":")
		path = ""

		for _, chunk := range pathChunks {
			// 4 is reserved for the colon and the ellipsis
			if len(path)+len(chunk) <= lineLimit-4 {
				path += chunk + ":"
			} else {
				path += "..."
				break
			}
		}
	}

	logger.UpdateSecondaryPrefix("setup")
	logger.Infof("PATH is now: %s", path)
	logger.ResetSecondaryPrefix()
}
