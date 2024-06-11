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

	testCase := test_cases.SingleLineOutputTestCase{
		Command:                    "type pwd",
		ExpectedPattern:            regexp.MustCompile(`^pwd is a( special)? shell builtin$`),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", (`pwd is a shell builtin`)),
		SuccessMessage:             "Received current working directory response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	var revertRenameOfPWD bool
	path, err := exec.LookPath("pwd")
	newPath := path + "Backup"
	// On MacOS, the OS doesn't allow renaming the `pwd` binary
	if err == nil && runtime.GOOS != "darwin" {
		revertRenameOfPWD = true
		// os.Rename is unable to complete this operation on some systems due to permission issues
		err = exec.Command("sudo", "mv", path, newPath).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Error renaming %q to %q: %v", path, newPath, err)
		}
	}

	testCase = test_cases.SingleLineOutputTestCase{
		Command:                    "pwd",
		ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`^%s$`, cwd)),
		ExpectedPatternExplanation: fmt.Sprintf("match %q", cwd+"\n"),
		SuccessMessage:             "Received current working directory response",
	}
	err = testCase.Run(shell, logger)

	if revertRenameOfPWD {
		err = exec.Command("sudo", "mv", newPath, path).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Error renaming %q to %q: %v", newPath, path, err)
		}
	}

	if err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
