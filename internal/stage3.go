package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
	b := shell_executable.NewShellExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	tries := random.RandomInt(3, 5)
	a := assertions.BufferAssertion{}
	truncatedStdErrBuf := shell_executable.NewTruncatedBuffer(b.GetStdErrBuffer())

	for i := 0; i < tries; i++ {
		command := "nonexistent" + strconv.Itoa(i)
		b.FeedStdin([]byte(command))
		expectedErrorMessage := fmt.Sprintf("%s: command not found", command)

		a.ExpectedValue = expectedErrorMessage
		if err := a.Run(&truncatedStdErrBuf, assertions.CoreTestInexact); err != nil {
			return err
		}

		truncatedStdErrBuf.UpdateOffsetToCurrentLength()
		logger.Debugf("Received message: %q", a.ActualValue)

		if strings.Contains(a.ActualValue, "\n") {
			lines := strings.Split(a.ActualValue, "\n")
			if len(lines) > 2 {
				a.ActualValue = lines[len(lines)-2]
			}
		}
		logger.Successf("Received error message: %q", a.ActualValue)

		if b.HasExited() {
			return fmt.Errorf("Program exited before all commands were sent")
		}
	}

	return nil
}
