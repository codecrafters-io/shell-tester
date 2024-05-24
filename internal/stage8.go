package internal

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/elf_executable"
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

	randomCode := elf_executable.GetRandomString()
	randomName := elf_executable.GetRandomName()
	err := custom_executable.CreateExecutable(randomCode, "my_exe")
	if err != nil {
		return err
	}

	// Add the current directory to PATH
	// (That is where the my_exe file is created)
	homeDir, _ := os.Getwd()
	path := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s", homeDir, path))

	command := []string{
		"./my_exe", randomName,
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
