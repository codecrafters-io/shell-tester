package internal

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	numberOfCommands := random.RandomInt(3, 6)

	if err := shell.Start(); err != nil {
		return err
	}

	for i := 0; i < numberOfCommands; i++ {
		command := "invalid_command_" + strconv.Itoa(i+1)

		if err := shell.AssertPrompt("$ "); err != nil {
			return err
		}

		if err := shell.SendCommand(command); err != nil {
			return err
		}

		if err := shell.AssertOutputMatchesRegex(regexp.MustCompile(fmt.Sprintf(`%s: (command )?not found\r\n`, command))); err != nil {
			return err
		}

		logger.Successf("✓ Received command not found message")
	}

	// At this stage the user might or might not have implemented a REPL to print the prompt again, so we won't test further

	return nil
}

// // 	logger := stageHarness.Logger
// // 	tries := random.RandomInt(3, 5)
// // 	a := assertions.BufferAssertion{}
// // 	truncatedStdErrBuf := shell_executable.NewTruncatedBuffer(b.GetStdErrBuffer())

// // 	for i := 0; i < tries; i++ {
// // 		command := "nonexistent" + strconv.Itoa(i)
// // 		b.FeedStdin([]byte(command))
// // 		expectedErrorMessage := fmt.Sprintf("%s: command not found", command)

// // 		a.ExpectedValue = expectedErrorMessage
// // 		if err := a.Run(&truncatedStdErrBuf, assertions.CoreTestInexact); err != nil {
// // 			return err
// // 		}

// // 		truncatedStdErrBuf.UpdateOffsetToCurrentLength()
// // 		logger.Debugf("Received message: %q", a.ActualValue)

// // 		if strings.Contains(a.ActualValue, "\n") {
// // 			lines := strings.Split(a.ActualValue, "\n")
// // 			if len(lines) > 2 {
// // 				a.ActualValue = lines[len(lines)-2]
// // 			}
// // 		}
// // 		logger.Successf("Received error message: %q", a.ActualValue)

// // 		if b.HasExited() {
// // 			return fmt.Errorf("Program exited before all commands were sent")
// // 		}
// // 	}

// // 	return nil
// // }

// func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
// 	return nil
// }
