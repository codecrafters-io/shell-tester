package internal

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRun(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomDir, err := GetRandomDirectory()
	if err != nil {
		return err
	}

	// Add randomDir to PATH (That is where the my_exe file is created)
	currentPath := os.Getenv("PATH")
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, currentPath))

	if err := shell.Start(); err != nil {
		return err
	}

	arch := runtime.GOARCH
	if arch == "arm64" {
		// The statically linked exe will not run on darwin arm64, just return nil
		return nil
	}

	randomCode := GetRandomString()
	randomName := GetRandomName()
	exePath := path.Join(randomDir, "my_exe")

	err = custom_executable.CreateExecutable(randomCode, exePath)
	if err != nil {
		return err
	}

	command := []string{
		"my_exe", randomName,
	}

	expectedResponseRegex := fmt.Sprintf("^Hello %s! The secret code is %s.$", randomName, randomCode)
	expectedResponse := fmt.Sprintf("Hello %s! The secret code is %s.", randomName, randomCode)
	testCase := test_cases.SingleLineOutputTestCase{
		Command:                    strings.Join(command, " "),
		ExpectedPattern:            regexp.MustCompile(expectedResponseRegex),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", expectedResponse),
		SuccessMessage:             "Received expected response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
