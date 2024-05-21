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
	args      []string
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
	b.args = args
	b.logger.Infof(b.getInitialLogLine())

	cmd := exec.Command(b.executable.Path)

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

func (b *ShellExecutable) LogOutput(output []byte) {
	b.programLogger.Plainln(string(output))
}

func (b *ShellExecutable) ReadBytesUntil(condition func([]byte) bool) ([]byte, error) {
	return b.ptyReader.ReadUntilCondition(condition)
}

// func (b *ShellExecutable) AssertOutputMatchesRegex(regexp *regexp.Regexp) error {
// 	regexMatchCondition := func(buf []byte) bool {
// 		return regexp.Match(StripANSI(buf))
// 	}

// 	actualValue, err := b.ptyReader.ReadUntilCondition(regexMatchCondition)
// 	if len(actualValue) > 0 {
// 		b.programLogger.Plainf("%s", string(StripANSI(actualValue)))
// 	}

// 	return nil
// }

func (b *ShellExecutable) AssertPrompt(prompt string) error {
	matchesPromptCondition := func(buf []byte) bool {
		return string(StripANSI(buf)) == prompt
	}

	actualValue, err := b.ptyReader.ReadUntilCondition(matchesPromptCondition)

	if err != nil {
		// If the user sent any output, let's print it before the error message.
		if len(actualValue) > 0 {
			b.LogOutput(StripANSI(actualValue))
		}

		return fmt.Errorf("Expected %q, got %q", prompt, string(actualValue))
	}

	extraOutput, extraOutputErr := b.ptyReader.ReadUntilTimeout(10 * time.Millisecond)
	fullOutput := append(actualValue, extraOutput...)

	// Whether the value matches our expecations or not, we print it
	b.LogOutput(StripANSI(fullOutput))

	// We failed to read extra output
	if extraOutputErr != nil {
		return fmt.Errorf("Error reading output: %v", extraOutputErr)
	}

	if len(extraOutput) > 0 {
		return fmt.Errorf("Found extra output after prompt: %q. (expected just %q)", string(extraOutput), prompt)
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

	reflectionCondition := func(buf []byte) bool {
		return string(buf) == expectedReflection
	}

	readBytes, err := b.ptyReader.ReadUntilCondition(reflectionCondition)
	if err != nil {
		return fmt.Errorf("Expected %q, but got %q", expectedReflection, string(readBytes))
	}

	return nil
}

func (b *ShellExecutable) getInitialLogLine() string {
	var log string

	if b.args == nil || len(b.args) == 0 {
		log = ("Running ./spawn_shell.sh")
	} else {
		log += "Running ./spawn_shell.sh"
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

// TODO: Use a library for this
func StripANSI(data []byte) []byte {
	// https://github.com/acarl005/stripansi/blob/master/stripansi.go
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	var re = regexp.MustCompile(ansi)

	return re.ReplaceAll(data, []byte(""))
}
