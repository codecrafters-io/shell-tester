package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testExit(stageHarness *test_case_harness.TestCaseHarness) error {
	b := shell_executable.NewShellExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	command := "nonexistent"
	b.FeedStdin([]byte(command))

	buffer, err := b.ReadBuffer("stderr")
	if err != nil {
		return err
	}

	errorMessage := string(buffer)

	if !strings.Contains(errorMessage, command+": command not found") {
		return fmt.Errorf("Expected error message to contain '%s: command not found', but got '%s'", command, errorMessage)
	}
	logger.Successf(strings.Split(errorMessage, "\n")[1])

	if b.HasExited() {
		return fmt.Errorf("Program exited before all commands were sent")
	}

	// // BUG: ReadBuffer's condition is breaking, because shell is only sending the output to stdout
	// // Not the prompt or any other output
	// b.FeedStdin([]byte("echo foo"))
	// buffer, err = b.ReadBuffer("stderr")
	// if err != nil {
	// 	return err
	// }
	// response = string(buffer)
	// fmt.Println("response", errorMessage)

	// Note: When we run exit, it exits the shell with the status of the last command executed, which is still 127 in this case.
	b.FeedStdin([]byte("exit 0"))
	result, err := b.Wait()
	if err != nil {
		return err
	}

	if result.ExitCode == -1 {
		return fmt.Errorf("Program did not exit after sending 'exit'")
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("Program did not exit with status 0 after sending 'exit'")
	}

	return nil
}
