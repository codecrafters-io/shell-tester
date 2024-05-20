package internal

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/creack/pty"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	cmd := exec.Command("bash")
	cmd.Env = []string{}
	// TODO: Find a way to take current environment bu sanitize ZSH-specific stuff
	cmd.Env = append(cmd.Env, "ZDOTDIR=/Users/rohitpaulk/experiments/codecrafters/testers/shell-tester/internal/test_helpers/zsh_config/")
	cmd.Env = append(cmd.Env, "BASH_SILENCE_DEPRECATION_WARNING=1")
	cmd.Env = append(cmd.Env, "PS1=$ ")
	cmd.Env = append(cmd.Env, "TERM=dumb")
	// cmd := exec.Command("sleep", "5")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	// Read prompt: How do we filter out ANSI sequences **after** the prompt is printed?
	//    Solution: Wait until condition, and then after condition is met - sleep 5ms and read extra. (ANSI filter the extra stuff, and then check it is empty)
	time.Sleep(100 * time.Millisecond)
	doRead(ptmx)

	// Why is \r\n not echo-ed back, but \n is?
	sendAndReadInput(ptmx, "missing")

	// TODO: Why does the first line contain % and spaces and \r \r
	// TODO: Why does the reflection contain \r\r\n instead of \r\n like with bash

	time.Sleep(100 * time.Millisecond)
	doRead(ptmx)

	sendAndReadInput(ptmx, "missing2")
	doRead(ptmx)

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

func doRead(ptmx *os.File) {
	buf := make([]byte, 1024)
	n, err := ptmx.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Read %d bytes: %q\n", n, string(buf[:n]))
}

func sendAndReadInput(ptmx *os.File, input string) {
	ptmx.Write([]byte(input + "\n"))

	// Make this deterministic
	time.Sleep(100 * time.Millisecond)

	expectedReflection := input + "\r\n"

	receivedBuf := make([]byte, len(expectedReflection))

	n, err := ptmx.Read(receivedBuf)
	if err != nil {
		panic(err)
	}

	if n != len(expectedReflection) {
		fmt.Printf("Expected to read %d bytes, but read %d\n", len(expectedReflection), n)
		panic("Failed to read input we wrote")
	}

	if string(receivedBuf) != expectedReflection {
		fmt.Printf("Expected to read %q, but read %q\n", expectedReflection, string(receivedBuf))
		panic("Failed to read input we wrote")
	}
}

func getTermios(file *os.File) (*syscall.Termios, error) {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&termios)))
	if errno != 0 {
		return nil, errno
	}
	return &termios, nil
}

func setTermios(file *os.File, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		return errno
	}
	return nil
}
