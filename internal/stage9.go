package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testpwd(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Error getting cwd: %v", err)
	}

	testCase := test_cases.SingleLineExactMatchTestCase{
		Command:                    "type pwd",
		FallbackPatterns:           []*regexp.Regexp{regexp.MustCompile(`^pwd is a( special)? shell builtin$`)},
		ExpectedPatternExplanation: `pwd is a shell builtin`,
		SuccessMessage:             "Received 'pwd is a shell builtin'",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	path, pwdNotFoundErr := exec.LookPath("pwd")
	newPath := path + "Backup"

	moveCommand := "mv"
	_, sudoNotFoundErr := exec.LookPath("sudo")
	if sudoNotFoundErr == nil {
		moveCommand = "sudo" + " " + moveCommand
	}

	shellExecutablePath, err := exec.LookPath("sh")
	if err != nil {
		panic(fmt.Sprintf("CodeCrafters internal error. Cannot get location of sh executable: %v", err))
	}

	// On macOS, the OS doesn't allow renaming the `pwd` binary
	if pwdNotFoundErr == nil && runtime.GOOS != "darwin" {
		// os.Rename is unable to complete this operation on some systems due to permission issues
		command := fmt.Sprintf("%s %s %s", moveCommand, path, newPath)
		cmd := exec.Command(shellExecutablePath, "-c", command)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Command failed: %s, Error renaming %q to %q: %v", cmd.String(), path, newPath, err.Error())
		}

		revertCommand := fmt.Sprintf("%s %s %s", moveCommand, newPath, path)
		revertCmd := exec.Command(shellExecutablePath, "-c", revertCommand)

		defer func(command *exec.Cmd) {
			err := command.Run()
			if err != nil {
				logger.Errorf("CodeCrafters internal error. Command failed: %s, Error renaming %q to %q: %v", command.String(), newPath, path, err)
			}
		}(revertCmd)
	}

	testCase = test_cases.SingleLineExactMatchTestCase{
		Command:        "pwd",
		ExpectedOutput: cwd,
		SuccessMessage: "Received current working directory response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
