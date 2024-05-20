package shell_executable

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
)

type ShellExecutable struct {
	executable    *executable.Executable
	logger        *logger.Logger
	programLogger *logger.Logger

	// Set after starting
	args []string
	pty  *os.File
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	b := &ShellExecutable{
		executable:    stageHarness.NewExecutable(),
		logger:        stageHarness.Logger,
		programLogger: logger.GetLogger(stageHarness.Logger.IsDebug, "[your_program] "),
	}

	// stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *ShellExecutable) Start(args ...string) error {
	b.args = args
	log := b.getInitialLogLine()
	b.logger.Infof(log)

	cmd := exec.Command(b.executable.Path)

	// TODO: Find a way to take current environment but still sanitize ZSH-specific stuff
	cmd.Env = []string{}
	cmd.Env = append(cmd.Env, "ZDOTDIR=/Users/rohitpaulk/experiments/codecrafters/testers/shell-tester/internal/test_helpers/zsh_config/")
	cmd.Env = append(cmd.Env, "BASH_SILENCE_DEPRECATION_WARNING=1")
	cmd.Env = append(cmd.Env, "PS1=$ ")
	cmd.Env = append(cmd.Env, "TERM=dumb")

	pty, err := ptylib.Start(cmd)
	if err != nil {
		panic(err)
	}

	b.pty = pty

	return nil
}

func (b *ShellExecutable) AssertPrompt(prompt string) error {
	ptyBuffer := NewFileBuffer(b.pty)

	shouldStopReadingBuffer := func(buf []byte) error {
		if string(StripANSI(buf)) == prompt {
			return nil
		} else {
			return fmt.Errorf("Prompt not found")
		}
	}

	if actualValue, err := ptyBuffer.ReadBuffer(shouldStopReadingBuffer); err != nil {
		// b.logger.Debugf("Read bytes: %q", actualValue)
		return fmt.Errorf("Expected %q, but got %q", prompt, string(actualValue))
	} else {
		// b.logger.Debugf("Read bytes: %q", actualValue)
		b.programLogger.Plainf("%s", string(StripANSI(actualValue)))
	}

	return nil
}

func (b *ShellExecutable) SendCommand(command string) error {
	b.logger.Infof("> %s", command)

	if err := b.writeAndReadReflection(command); err != nil {
		return err
	}

	return nil
}

func (b *ShellExecutable) writeAndReadReflection(command string) error {
	b.pty.Write([]byte(command + "\n"))

	expectedReflection := command + "\r\n"
	readBytes := []byte{}

	for len(readBytes) < len(expectedReflection) {
		singleByteBuf := make([]byte, 1)
		n, err := b.pty.Read(singleByteBuf)
		if err != nil {
			if strings.Contains(err.Error(), "resource temporarily unavailable") {
				continue
			}

			return err
		}

		if n != 1 {
			return fmt.Errorf("Expected to read %d bytes, but read %d", len(expectedReflection), n)
		}

		readBytes = append(readBytes, singleByteBuf...)
	}

	if string(readBytes) != expectedReflection {
		return fmt.Errorf("Expected to read %q, but read %q", expectedReflection, string(readBytes))
	}

	return nil
}

// func (b *ShellExecutable) FeedStdin(command []byte) error {
// 	commandWithEnter := append(command, []byte("\n")...)
// 	return b.feedStdin(commandWithEnter)
// }

func (b *ShellExecutable) getInitialLogLine() string {
	var log string

	if b.args == nil || len(b.args) == 0 {
		log = ("$ ./spawn_shell.sh")
	} else {
		log += "$ ./spawn_shell.sh"
		for _, arg := range b.args {
			if strings.Contains(arg, " ") {
				log += " \"" + arg + "\""
			} else {
				log += " " + arg
			}
		}
	}

	return log
}

// func (b *ShellExecutable) HasExited() bool {
// 	return b.executable.HasExited()
// }

// func (b *ShellExecutable) Kill() error {
// 	b.logger.Debugf("Terminating program")
// 	if err := b.executable.Kill(); err != nil {
// 		b.logger.Debugf("Error terminating program: '%v'", err)
// 		return err
// 	}

// 	b.logger.Debugf("Program terminated successfully")
// 	return nil // When does this happen?
// }

// func (b *ShellExecutable) feedStdin(command []byte) error {
// 	n, err := b.executable.StdinPipe.Write(command)
// 	b.logger.Debugf("Wrote %d bytes to stdin", n)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (b *ShellExecutable) GetStdErrBuffer() *bytes.Buffer {
// 	return b.executable.StderrBuffer
// }

// func (b *ShellExecutable) GetStdOutBuffer() *bytes.Buffer {
// 	return b.executable.StdoutBuffer
// }

// func (b *ShellExecutable) Wait() (executable.ExecutableResult, error) {
// 	return b.executable.Wait()
// }
