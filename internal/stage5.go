package internal

// import "github.com/codecrafters-io/tester-utils/test_case_harness"

// // import (
// // 	"fmt"

// // 	"github.com/codecrafters-io/shell-tester/internal/assertions"
// // 	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
// // 	"github.com/codecrafters-io/tester-utils/test_case_harness"
// // )

// // func testEcho(stageHarness *test_case_harness.TestCaseHarness) error {
// // 	b := shell_executable.NewShellExecutable(stageHarness)
// // 	if err := b.Run(); err != nil {
// // 		return err
// // 	}

// // 	logger := stageHarness.Logger
// // 	message := "Hello World!"
// // 	command := fmt.Sprintf("echo %s", message)
// // 	b.FeedStdin([]byte(command))

// // 	a := assertions.BufferAssertion{ExpectedValue: message}
// // 	truncatedStdErrBuf := shell_executable.NewTruncatedBuffer(b.GetStdOutBuffer())
// // 	if err := a.Run(&truncatedStdErrBuf, assertions.CoreTestExact); err != nil {
// // 		return err
// // 	}
// // 	logger.Successf("Received message: %q", a.ActualValue)

// // 	if b.HasExited() {
// // 		return fmt.Errorf("Program exited before all commands were sent")
// // 	}

// // 	return nil
// // }

// func testEcho(stageHarness *test_case_harness.TestCaseHarness) error {
// 	return nil
// }
