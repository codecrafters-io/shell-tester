package internal

// import "github.com/codecrafters-io/tester-utils/test_case_harness"

// // import (
// // 	"fmt"
// // 	"strings"

// // 	"github.com/codecrafters-io/shell-tester/internal/assertions"
// // 	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
// // 	"github.com/codecrafters-io/tester-utils/test_case_harness"
// // )

// // func testExit(stageHarness *test_case_harness.TestCaseHarness) error {
// // 	b := shell_executable.NewShellExecutable(stageHarness)
// // 	if err := b.Run(); err != nil {
// // 		return err
// // 	}

// // 	logger := stageHarness.Logger
// // 	command := "nonexistent"
// // 	expectedErrorMessage := fmt.Sprintf("%s: command not found", command)
// // 	b.FeedStdin([]byte(command))

// // 	a := assertions.BufferAssertion{ExpectedValue: expectedErrorMessage}
// // 	truncatedStdErrBuf := shell_executable.NewTruncatedBuffer(b.GetStdErrBuffer())
// // 	if err := a.Run(&truncatedStdErrBuf, assertions.CoreTestInexact); err != nil {
// // 		return err
// // 	}
// // 	logger.Debugf("Received message: %q", a.ActualValue)

// // 	if strings.Contains(a.ActualValue, "\n") {
// // 		lines := strings.Split(a.ActualValue, "\n")
// // 		if len(lines) > 2 {
// // 			a.ActualValue = lines[len(lines)-2]
// // 		}
// // 	}

// // 	logger.Successf("Received error message: %q", a.ActualValue)

// // 	if b.HasExited() {
// // 		return fmt.Errorf("Program exited before all commands were sent")
// // 	}

// // 	b.FeedStdin([]byte("exit 0"))
// // 	result, err := b.Wait()
// // 	if err != nil {
// // 		return err
// // 	}

// // 	if result.ExitCode == -1 {
// // 		return fmt.Errorf("Program did not exit after sending 'exit'")
// // 	}
// // 	if result.ExitCode != 0 {
// // 		return fmt.Errorf("Expected exit code 0, but got %d", result.ExitCode)
// // 	}

// // 	return nil
// // }

// func testExit(stageHarness *test_case_harness.TestCaseHarness) error {
// 	return nil
// }
