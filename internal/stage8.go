package internal

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRun(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	randomCode := GetRandomString()
	randomName := GetRandomName()
	randomDir, err := GetRandomDirectory()
	if err != nil {
		return err
	}
	exePath := path.Join(randomDir, "my_exe")

	err = custom_executable.CreateExecutable(randomCode, exePath)
	if err != nil {
		return err
	}

	// Add the current directory to PATH
	// (That is where the my_exe file is created)
	path := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, path))

	command := []string{
		exePath, randomName,
	}

	expectedResponse := fmt.Sprintf("Hello %s! The secret code is %s.", randomName, randomCode)

	testCase := test_cases.RegexTestCase{
		Command:                    strings.Join(command, " "),
		ExpectedPattern:            regexp.MustCompile(expectedResponse),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", expectedResponse),
		SuccessMessage:             "Received expected response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
