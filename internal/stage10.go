package internal

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	directory, err := getRandomDirectory()
	if err != nil {
		return err
	}

	testCase := test_cases.CDAndPWDTestCase{Directory: directory, Response: directory}
	err = testCase.Run(shell, logger)
	if err != nil {
		return err
	}

	directory = "/non-existing-directory"
	command := fmt.Sprintf("cd %s", directory)

	failureTestCase := test_cases.RegexTestCase{
		Command:                    command,
		ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`(can't cd to %s|%s: No such file or directory)\r\n`, directory, directory)),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", fmt.Sprintf(`%s: No such file or directory\r\n`, directory)),
		SuccessMessage:             "Received error message",
	}

	if err := failureTestCase.Run(shell, logger); err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}

// getRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>/<random-word>/<random-word>`
func getRandomDirectory() (string, error) {
	randomDir := path.Join("/tmp", random.RandomWord(), random.RandomWord(), random.RandomWord())
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
	}
	return randomDir, nil
}
