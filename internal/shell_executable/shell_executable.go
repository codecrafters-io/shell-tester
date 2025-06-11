package shell_executable

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/condition_reader"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
	"go.chromium.org/luci/common/system/environ"
)

// ErrProgramExited is returned when the program exits
var ErrProgramExited = errors.New("Program exited")

type ShellExecutable struct {
	executable    *executable.Executable
	stageLogger   *logger.Logger
	programLogger *logger.Logger

	// env is set to os.Environ() by default, but individual values can be overridden with Setenv
	env environ.Env

	// Set after starting
	cmd       *exec.Cmd
	pty       *os.File
	ptyReader condition_reader.ConditionReader
	vt        *virtual_terminal.VirtualTerminal
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	b := &ShellExecutable{
		executable:    stageHarness.NewExecutable(),
		stageLogger:   stageHarness.Logger,
		programLogger: logger.GetLogger(stageHarness.Logger.IsDebug, "[your-program] "),
	}

	b.env = environ.New(os.Environ())

	// TODO: Kill pty process?
	// stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *ShellExecutable) Setenv(key, value string) {
	b.env.Set(key, value)
}

func (b *ShellExecutable) AddToPath(dir string) {
	b.env.Set("PATH", fmt.Sprintf("%s:%s", dir, b.env.Get("PATH")))
}

func (b *ShellExecutable) GetPath() string {
	return b.env.Get("PATH")
}

func (b *ShellExecutable) Start(args ...string) error {
	b.stageLogger.Infof("%s", b.getInitialLogLine(args...))

	b.Setenv("PS1", utils.PROMPT)
	// b.Setenv("TERM", "dumb") // test_all_success works without this too, do we need it?

	cmd := exec.Command(b.executable.Path, args...)
	cmd.Env = b.env.Sorted()

	b.cmd = cmd
	b.vt = virtual_terminal.NewStandardVT()

	winsize := &ptylib.Winsize{
		Rows: uint16(b.vt.GetRowCount()),
		Cols: uint16(b.vt.GetColumnCount()),
	}
	pty, err := ptylib.StartWithSize(cmd, winsize)
	if err != nil {
		return fmt.Errorf("Failed to execute %s: %v", b.executable.Path, err)
	}

	b.pty = pty
	b.ptyReader = condition_reader.NewConditionReader(io.TeeReader(b.pty, b.vt))

	return nil
}

func (b *ShellExecutable) GetScreenState() [][]string {
	return b.vt.GetScreenState()
}

func (b *ShellExecutable) LogOutput(output []byte) {
	b.programLogger.Plainf("%s", string(output))
}

func (b *ShellExecutable) VTBellChannel() chan bool {
	return b.vt.BellChannel()
}

func (b *ShellExecutable) ReadUntilConditionOrTimeout(condition func() bool, timeout time.Duration) error {
	err := b.ptyReader.ReadUntilConditionOrTimeout(condition, timeout)
	if err != nil {
		return wrapReaderError(err)
	}

	return nil
}

func (b *ShellExecutable) SendCommand(command string) error {
	if err := b.SendCommandRaw(command + "\n"); err != nil {
		return err
	}

	return nil
}

func (b *ShellExecutable) SendCommandRaw(command string) error {
	if _, err := b.pty.Write([]byte(command)); err != nil {
		return err
	}

	return nil
}

func (b *ShellExecutable) WaitForTermination() (hasTerminated bool, exitCode int) {
	if b.cmd == nil {
		panic("CodeCrafters internal error: WaitForTermination called before command was run")
	}

	waitCompleted := make(chan bool)

	go func() {
		b.cmd.Wait()
		waitCompleted <- true
	}()

	select {
	case <-waitCompleted:
		rawExitCode := b.cmd.ProcessState.ExitCode()

		if rawExitCode == -1 {
			// We can get isTerminated as false if the program is terminated by SIGKILL too, but that seems unlikely here
			return false, 0
		} else {
			return true, rawExitCode
		}
	case <-time.After(2 * time.Second):
		return false, 0
	}
}

func (b *ShellExecutable) ExitCode() int {
	// Calling WaitForTermination multiple times is okay, Wait() would error out, but we will get the exit code
	exited, exitCode := b.WaitForTermination()
	if !exited {
		// fmt.Println("Process has not exited yet.")
		return -1
	}
	return exitCode
}

func (b *ShellExecutable) getInitialLogLine(args ...string) string {
	var log string

	if len(args) == 0 {
		log = fmt.Sprintf("Running ./%s", path.Base(b.executable.Path))
	} else {
		log += fmt.Sprintf("Running ./%s", path.Base(b.executable.Path))
		for _, arg := range args {
			if strings.Contains(arg, " ") {
				log += " \"" + arg + "\""
			} else {
				log += " " + arg
			}
		}
	}

	return log
}

func wrapReaderError(readerErr error) error {
	// Linux returns syscall.EIO when the process is killed, macOS returns io.EOF
	if errors.Is(readerErr, io.EOF) || errors.Is(readerErr, syscall.EIO) {
		return ErrProgramExited
	}

	return readerErr
}

func (b *ShellExecutable) GetLogger() *logger.Logger {
	return b.stageLogger
}
