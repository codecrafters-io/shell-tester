package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRun(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := getRandomDirectory()
	if err != nil {
		return err
	}

	// Add randomDir to PATH (That is where the my_exe file is created)
	currentPath := os.Getenv("PATH")
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, currentPath))

	if err := shell.Start(); err != nil {
		return err
	}

	// First prompt assertion
	if err := asserter.Assert(); err != nil {
		return err
	}

	randomCode := getRandomString()
	randomName := getRandomName()
	exePath := path.Join(randomDir, "my_exe")

	err = custom_executable.CreateExecutable(randomCode, exePath)
	if err != nil {
		return err
	}

	command := []string{
		"my_exe", randomName,
	}

	testCase := test_cases.CommandResponseTestCase{
		Command:          strings.Join(command, " "),
		ExpectedOutput:   fmt.Sprintf("Hello %s! The secret code is %s", randomName, randomCode),
		FallbackPatterns: nil,
		SuccessMessage:   "Received expected response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	asserter.LogRemainingOutput()
	return nil
}
