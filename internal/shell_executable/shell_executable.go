package shell_executable

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/condition_reader"
	"github.com/codecrafters-io/shell-tester/internal/screen_state"
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

// ErrMemoryLimitExceeded is returned when a process exceeds its memory limit
var ErrMemoryLimitExceeded = errors.New("process exceeded memory limit")

const defaultMemoryLimitBytes = 2 * 1024 * 1024 * 1024 // 2GB

type ShellExecutable struct {
	// MemoryLimitInBytes sets the maximum memory the process can use (Linux only).
	// If exceeded, the process will be killed and an error will be returned.
	// Defaults to 2GB. Set to 0 to disable memory limiting.
	MemoryLimitInBytes int64

	memoryMonitor *memoryMonitor // Monitors process memory usage and kills if limit exceeded
	oomKilled     bool           // set when monitor is stopped after process exits

	executable    *executable.Executable
	stageLogger   *logger.Logger
	programLogger *logger.Logger

	// env is set to os.Environ() by default, but individual values can be overridden with Setenv
	env environ.Env

	workingDir string

	// Set after starting
	cmd       *exec.Cmd
	pty       *os.File
	ptyReader condition_reader.ConditionReader
	vt        *virtual_terminal.VirtualTerminal
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	b := &ShellExecutable{
		executable:         stageHarness.NewExecutable(),
		stageLogger:        stageHarness.Logger,
		programLogger:      logger.GetLogger(stageHarness.Logger.IsDebug, "[your-program] "),
		MemoryLimitInBytes: defaultMemoryLimitBytes,
	}

	b.env = environ.New(os.Environ())

	// TODO: Kill pty process?
	// stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

// NewShellExecutableForTest creates a ShellExecutable that runs the executable at path.
// Used for testing (e.g. memory limit tests). The logger can be nil; a quiet logger will be used.
func NewShellExecutableForTest(path string, stageLogger *logger.Logger) *ShellExecutable {
	if stageLogger == nil {
		stageLogger = logger.GetQuietLogger("")
	}
	b := &ShellExecutable{
		executable:         executable.NewExecutable(path),
		stageLogger:        stageLogger,
		programLogger:      logger.GetLogger(stageLogger.IsDebug, "[your-program] "),
		MemoryLimitInBytes: defaultMemoryLimitBytes,
	}
	b.env = environ.New(os.Environ())
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

func (b *ShellExecutable) SetWorkingDirectory(workingDirPath string) {
	b.workingDir = workingDirPath
}

func (b *ShellExecutable) Start(args ...string) error {
	b.stageLogger.Infof("%s", b.getInitialLogLine(args...))

	b.Setenv("PS1", utils.PROMPT)
	// b.Setenv("TERM", "dumb") // test_all_success works without this too, do we need it?

	absolutePath, err := filepath.Abs(b.executable.Path)
	if err != nil {
		return err
	}

	cmd := exec.Command(absolutePath, args...)
	// If workingDir is empty, it is set as cwd() by exec library
	cmd.Dir = b.workingDir
	cmd.Env = b.env.Sorted()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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

	// Start memory monitoring for RSS-based memory limiting (Linux only, no-op on other platforms)
	b.memoryMonitor = newMemoryMonitor(b.MemoryLimitInBytes)
	b.memoryMonitor.start(cmd.Process.Pid)

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
	if err := b.SendText(command + "\n"); err != nil {
		return err
	}

	return nil
}

func (b *ShellExecutable) SendText(text string) error {
	if _, err := b.pty.Write([]byte(text)); err != nil {
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
		// Stop memory monitor and cache OOM result before clearing
		if b.memoryMonitor != nil {
			b.oomKilled = b.memoryMonitor.wasOOMKilled()
			b.memoryMonitor.stop()
			b.memoryMonitor = nil
		}

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

// formatBytesHumanReadable formats bytes as a human-readable string (e.g., "50 MB", "2 GB")
func formatBytesHumanReadable(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%d GB", bytes/GB)
	case bytes >= MB:
		return fmt.Sprintf("%d MB", bytes/MB)
	case bytes >= KB:
		return fmt.Sprintf("%d KB", bytes/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// WasOOMKilled returns true if the process was killed due to exceeding the memory limit.
// Only meaningful after the process has terminated (e.g. after WaitForTermination returns true).
func (b *ShellExecutable) WasOOMKilled() bool {
	return b.oomKilled
}

// MemoryLimitExceededError returns an error describing that the process exceeded its memory limit,
// or nil if the process was not OOM killed. Use after WaitForTermination returns true.
func (b *ShellExecutable) MemoryLimitExceededError() error {
	if !b.oomKilled {
		return nil
	}
	return fmt.Errorf("process exceeded memory limit (%s): %w", formatBytesHumanReadable(b.MemoryLimitInBytes), ErrMemoryLimitExceeded)
}

func (b *ShellExecutable) GetLogger() *logger.Logger {
	return b.stageLogger
}
