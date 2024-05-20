package internal

import (
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	shell := shell_executable.NewShellExecutable(stageHarness)
	if err := shell.Start(); err != nil {
		return err
	}

	if err := shell.AssertPrompt("$ "); err != nil {
		return err
	}

	if err := shell.SendCommand("missing"); err != nil {
		return err
	}

	if err := shell.AssertPrompt("bash: blah: command not found"); err != nil {
		return err
	}

	// doRead(ptmx)

	// sendAndReadInput(ptmx, "missing2")
	// doRead(ptmx)

	return nil

	// logger := stageHarness.Logger
	// command := "nonexistent"
	// expectedErrorMessage := fmt.Sprintf("%s: command not found", command)

	// a := assertions.BufferAssertion{ExpectedValue: expectedErrorMessage}
	// stdBuffer := shell_executable.NewFileBuffer(ptmx)
	// stdBuffer.FeedStdin([]byte(command))

	// if err := a.Run(&stdBuffer, assertions.CoreTestInexact); err != nil {
	// 	return err
	// }
	// logger.Debugf("Received message: %q", a.ActualValue)

	// if strings.Contains(a.ActualValue, "\n") {
	// 	lines := strings.Split(a.ActualValue, "\n")
	// 	if len(lines) >= 2 {
	// 		a.ActualValue = lines[len(lines)-2]
	// 	}
	// }

	// logger.Successf("Received error message: %q", a.ActualValue)
	// return nil
}
