package internal

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/creack/pty"
)

func testPrompt(stageHarness *test_case_harness.TestCaseHarness) error {
	os.Setenv("PS1", "$ ")
	os.Setenv("BASH_SILENCE_DEPRECATION_WARNING", "1")

	cmd := exec.Command("bash", "--norc", "-i")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	logger := stageHarness.Logger

	expectedPrompt := "$"

	a := assertions.BufferAssertion{ExpectedValue: expectedPrompt}
	stdBuffer := shell_executable.NewFileBuffer(ptmx)

	if err := a.Run(&stdBuffer, assertions.CoreTestInexact); err != nil {
		return err
	}
	logger.Successf("Received prompt: %q", a.ActualValue)

	if cmd.ProcessState != nil {
		return fmt.Errorf("Expected shell to be running, but it has exited")
	}
	logger.Successf("Shell is still running")

	cmd.Process.Kill()
	cmd.Wait()
	fmt.Println("cmd.ProcessState: ", cmd.ProcessState)

	return nil
}
