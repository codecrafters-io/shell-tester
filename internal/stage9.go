package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testpwd(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Error getting cwd: %v", err)
	}

	testCase := test_cases.CommandResponseTestCase{
		Command:          "type pwd",
		ExpectedOutput:   `pwd is a shell builtin`,
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile(`^pwd is a( special)? shell builtin$`)},
		SuccessMessage:   "Received 'pwd is a shell builtin'",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	path, pwdNotFoundErr := exec.LookPath("pwd")
	newPath := path + "Backup"

	moveCommand := "mv"
	_, sudoNotFoundErr := exec.LookPath("sudo")
	if sudoNotFoundErr == nil {
		moveCommand = "sudo" + " " + moveCommand
	}
	// On macOS, the OS doesn't allow renaming the `pwd` binary
	if pwdNotFoundErr == nil && runtime.GOOS != "darwin" {
		// os.Rename is unable to complete this operation on some systems due to permission issues
		command := fmt.Sprintf("%s %s %s", moveCommand, path, newPath)
		err = exec.Command("sh", "-c", command).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters internal error. Error renaming %q to %q: %v", path, newPath, err)
		}

		revertCommand := fmt.Sprintf("%s %s %s", moveCommand, newPath, path)

		defer func(command *exec.Cmd) {
			err := command.Run()
			if err != nil {
				logger.Errorf("CodeCrafters internal error. Error renaming %q to %q: %v", newPath, path, err)
			}
		}(exec.Command("sh", "-c", revertCommand))
	}

	testCase = test_cases.CommandResponseTestCase{
		Command:          "pwd",
		ExpectedOutput:   cwd,
		FallbackPatterns: nil,
		SuccessMessage:   "Received current working directory response",
	}
	if err := testCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
