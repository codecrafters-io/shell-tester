package shell_executable

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/condition_reader"
	"github.com/codecrafters-io/shell-tester/internal/screen_state"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	virtual_terminal "github.com/codecrafters-io/shell-tester/internal/vt"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/executable/stdio_handler"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
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
	ptyReader condition_reader.ConditionReader
	vt        *virtual_terminal.VirtualTerminal
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	executable := executable.NewVerboseExecutable(
		stageHarness.Executable.Path,
		func(s string) {},
	)

	b := &ShellExecutable{
		executable:    executable,
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

	b.vt = virtual_terminal.NewStandardVT()

	// Supply environment variables beore starting the process
	b.executable.StdioHandler = &stdio_handler.SinglePtyStdioHandler{
		Width:  uint(b.vt.GetColumnCount()),
		Height: uint(b.vt.GetRowCount()),
	}

	b.executable.Env = b.env

	if err := b.executable.Start(args...); err != nil {
		return err
	}

	b.ptyReader = condition_reader.NewConditionReader(
		io.TeeReader(b.executable.GetStdoutStreamReader(), b.vt),
	)

	return nil
}

func (b *ShellExecutable) GetScreenState() screen_state.ScreenState {
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
	stdinDevice := b.executable.GetStdinWriter()

	if _, err := stdinDevice.Write([]byte(command)); err != nil {
		return err
	}

	return nil
}

func (b *ShellExecutable) Wait() (executable.ExecutableResult, error) {
	return b.executable.Wait()
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
