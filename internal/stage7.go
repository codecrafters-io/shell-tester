package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func getPath(executable string) string {
	path, err := exec.LookPath(executable)
	if err != nil {
		return fmt.Sprintf(`%s[:]? not found`, executable)
	} else {
		return path
	}
}

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	// Add the current directory to PATH (That is where the my_exe file is created)
	//
	// TODO: Remove this since it mutates path for ALL stages! Use shell.Setenv() instead.
	//       We'll need to change the test to not use exec.LookPath for my_exe at least.
	//       Also, always create my_exe in a RANDOM directory, not the current directory.
	homeDir, _ := os.Getwd()
	path := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s", homeDir, path))

	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	err := custom_executable.CreateExecutable(GetRandomString(), "my_exe")
	if err != nil {
		return err
	}

	if err := shell.Start(); err != nil {
		return err
	}

	availableExecutables := []string{"cat", "cp", "mkdir", "my_exe"}

	for _, executable := range availableExecutables {
		command := fmt.Sprintf("type %s", executable)
		actualPath := getPath(executable)
		expectedPattern := fmt.Sprintf(`^(%s is )?%s\r\n`, executable, actualPath)
		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(expectedPattern),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", fmt.Sprintf(`%s is %s`, executable, actualPath)),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	nonAvailableExecutables := []string{"nonexistent"}

	for _, executable := range nonAvailableExecutables {
		command := fmt.Sprintf("type %s", executable)
		actualPath := getPath(executable)
		expectedPattern := fmt.Sprintf(`^(bash: type: )?%s\r\n`, actualPath)
		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(expectedPattern),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", fmt.Sprintf(`%s: not found`, executable)),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}
