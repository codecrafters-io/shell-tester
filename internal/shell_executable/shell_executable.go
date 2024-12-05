package shell_executable

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/condition_reader"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
	"go.chromium.org/luci/common/system/environ"
)

// ErrConditionNotMet is re-exported from condition_reader for convenience
var ErrConditionNotMet = condition_reader.ErrConditionNotMet

// ErrProgramExited is returned when the program exits
var ErrProgramExited = errors.New("Program exited")

type ShellExecutable struct {
	executable    *executable.Executable
	logger        *logger.Logger
	programLogger *logger.Logger

	// env is set to os.Environ() by default, but individual values can be overridden with Setenv
	env environ.Env

	// Set after starting
	cmd       *exec.Cmd
	pty       *os.File
	ptyReader condition_reader.ConditionReader
	vt        *VirtualTerminal
	termIO    *TermIO
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	b := &ShellExecutable{
		executable:    stageHarness.NewExecutable(),
		logger:        stageHarness.Logger,
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

func (b *ShellExecutable) Start(args ...string) error {
	b.logger.Infof(b.getInitialLogLine(args...))

	b.Setenv("PS1", "$ ")
	// b.Setenv("TERM", "dumb") // test_all_success works without this too, do we need it?

	cmd := exec.Command(b.executable.Path, args...)
	cmd.Env = b.env.Sorted()

	pty, err := ptylib.Start(cmd)
	if err != nil {
		return fmt.Errorf("Failed to execute %s: %v", b.executable.Path, err)
	}

	b.cmd = cmd
	b.pty = pty
	b.vt = NewStandardVT()
	b.termIO = NewTermIO(b.vt, b.pty)
	b.ptyReader = condition_reader.NewConditionReader(b.termIO)
	// defer b.vt.Close() // ToDo ??

	return nil
}

func (b *ShellExecutable) GetScreenState(retainColors bool) [][]string {
	return b.vt.GetScreenState(retainColors)
}

func (b *ShellExecutable) GetScreenStateSingleRow(row int, retainColors bool) []string {
	return b.vt.GetRow(row, retainColors)
}

func (b *ShellExecutable) GetScreenStateForLogging(retainColors bool) string {
	fullScreenState := b.GetScreenState(retainColors)
	screenStateString := ""
	for _, row := range fullScreenState {
		var filteredRow []string
		for _, cell := range row {
			if cell != "." {
				filteredRow = append(filteredRow, cell)
			}
		}
		if len(filteredRow) == 0 {
			continue
		}
		screenStateString += strings.Join(filteredRow, "")
		screenStateString += "\n"
	}
	return screenStateString
}

func (b *ShellExecutable) GetRowsTillEndForLogging(startingRow int, retainColors bool) string {
	fullScreenState := b.vt.GetRowsTillEnd(startingRow, retainColors)
	screenStateString := ""
	for _, row := range fullScreenState {
		var filteredRow []string
		for _, cell := range row {
			if cell != "." {
				filteredRow = append(filteredRow, cell)
			}
		}
		if len(filteredRow) == 0 {
			continue
		}
		screenStateString += strings.Join(filteredRow, "")
		screenStateString += "\n"
	}
	return screenStateString
}

func (b *ShellExecutable) GetScreenStateSingleRowForLogging(row int, retainColors bool) string {
	screenStateSingleRow := b.GetScreenStateSingleRow(row, retainColors)

	screenStateString := ""
	var filteredRow []string
	for _, cell := range screenStateSingleRow {
		if cell != "." {
			filteredRow = append(filteredRow, cell)
		}
	}
	if len(filteredRow) == 0 {
		return ""
	}
	screenStateString += strings.Join(filteredRow, "")
	screenStateString += "\n"

	return screenStateString
}

// TODO: Do tests cases _need_ to decide when to log output and when to not? Can we just always log from within ReadBytes...?

func (b *ShellExecutable) LogOutput(output []byte) {
	b.programLogger.Plainln(string(output))
}

func (b *ShellExecutable) ReadBytesUntil(condition func([]byte) bool) ([]byte, error) {
	readBytes, err := b.ptyReader.ReadUntilCondition(condition)
	if err != nil {
		return readBytes, wrapReaderError(err)
	}

	return readBytes, nil
}

func (b *ShellExecutable) ReadBytesUntilTimeout(timeout time.Duration) ([]byte, error) {
	readBytes, err := b.ptyReader.ReadUntilTimeout(timeout)
	if err != nil {
		return readBytes, wrapReaderError(err)
	}

	return readBytes, nil
}

func (b *ShellExecutable) SendCommand(command string) error {
	b.logger.Infof("> %s", command)

	if err := b.writeAndReadReflection(command); err != nil {
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

func (b *ShellExecutable) writeAndReadReflection(command string) error {
	b.pty.Write([]byte(command + "\n"))

	expectedReflection := command + "\r\n"
	var readBytes []byte

	reflectionCondition := func(buf []byte) bool {
		return string(buf) == expectedReflection
	}

	readBytes, err := b.ptyReader.ReadUntilCondition(reflectionCondition)
	if err != nil {
		return fmt.Errorf("CodeCrafters internal error. Expected %q when writing to pty, but got %q", expectedReflection, string(readBytes))
	}

	return nil
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

func StripANSI(data []byte) []byte {
	// https://github.com/acarl005/stripansi/blob/master/stripansi.go
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	var re = regexp.MustCompile(ansi)

	return re.ReplaceAll(data, []byte(""))
}

func wrapReaderError(readerErr error) error {
	// Linux returns syscall.EIO when the process is killed, macOS returns io.EOF
	if errors.Is(readerErr, io.EOF) || errors.Is(readerErr, syscall.EIO) {
		return ErrProgramExited
	}

	return readerErr
}
