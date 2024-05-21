package shell_executable

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/condition_reader"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
)

// ErrConditionNotMet is re-exported from condition_reader for convenience
var ErrConditionNotMet = condition_reader.ErrConditionNotMet

type ShellExecutable struct {
	executable    *executable.Executable
	logger        *logger.Logger
	programLogger *logger.Logger

	// Set after starting
	pty       *os.File
	ptyReader condition_reader.ConditionReader
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	b := &ShellExecutable{
		executable:    stageHarness.NewExecutable(),
		logger:        stageHarness.Logger,
		programLogger: logger.GetLogger(stageHarness.Logger.IsDebug, "[your-program] "),
	}

	// TODO: Kill pty process?
	// stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *ShellExecutable) Start(args ...string) error {
	b.logger.Infof(b.getInitialLogLine(args...))

	cmd := exec.Command(b.executable.Path, args...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PS1=$ ")
	cmd.Env = append(cmd.Env, "TERM=dumb") // test_all_success works without this too, do we need it?

	pty, err := ptylib.Start(cmd)
	if err != nil {
		return fmt.Errorf("Failed to execute %s: %v", b.executable.Path, err)
	}

	b.pty = pty
	b.ptyReader = condition_reader.NewConditionReader(b.pty)

	return nil
}

// TODO: Do tests cases _need_ to decide when to log output and when to not? Can we just always log from within ReadBytes...?
func (b *ShellExecutable) LogOutput(output []byte) {
	b.programLogger.Plainln(string(output))
}

func (b *ShellExecutable) ReadBytesUntil(condition func([]byte) bool) ([]byte, error) {
	return b.ptyReader.ReadUntilCondition(condition)
}

func (b *ShellExecutable) ReadBytesUntilTimeout(timeout time.Duration) ([]byte, error) {
	return b.ptyReader.ReadUntilTimeout(timeout)
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
		log = ("Running ./spawn_shell.sh")
	} else {
		log += "Running ./spawn_shell.sh"
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
