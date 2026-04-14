package test_cases

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/google/shlex"
	"github.com/shirou/gopsutil/v3/process"
)

// BackgroundCommandResponseTestCase launches the given command with an & symbol
// Launching it to the background
// It will assert that the job number is the expected one in the output
// It asserts the next prompt immediately
type BackgroundCommandResponseTestCase struct {
	Command           string
	ExpectedJobNumber int
	SuccessMessage    string
}

func (t *BackgroundCommandResponseTestCase) Run(asserter *logged_shell_asserter.LoggedShellAsserter, shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	commandToSend := fmt.Sprintf("%s &", t.Command)

	if err := shell.SendCommand(commandToSend); err != nil {
		return fmt.Errorf("Error sending command to shell: %v", err)
	}

	commandReflection := fmt.Sprintf("$ %s", commandToSend)
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: commandReflection,
	})

	// Assert command reflection first
	if err := asserter.AssertWithoutPrompt(); err != nil {
		return err
	}

	outputFormatRegex := regexp.MustCompile(
		fmt.Sprintf(`^\[%d\] (\d+)$`, t.ExpectedJobNumber),
	)

	// Assert the output format first
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: fmt.Sprintf("[%d] <PID>", t.ExpectedJobNumber),
		FallbackPatterns: []*regexp.Regexp{
			outputFormatRegex,
		},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// Extract the PID from the output format and check if that PID is the shell's child
	outputLineIdx := asserter.GetLastLoggedRowIndex()
	outputLine := shell.GetScreenState().GetRow(outputLineIdx)

	matches := outputFormatRegex.FindStringSubmatch(outputLine.String())

	// This will never trigger: This was already asserted
	if matches == nil {
		panic("Codecrafters Internal Error: Could not match PID from background command output")
	}

	receivedPid, err := strconv.Atoi(matches[1])

	// This will never trigger: This was already assserted
	if err != nil {
		panic("Codecrafters Internal Error: Could not parse PID from background command output")
	}

	childPids := shell.GetAllDescendentsPids()

	if !slices.Contains(childPids, receivedPid) {
		return fmt.Errorf("Could not find process with PID %d", receivedPid)
	}

	logger.Successf("✓ Found process with PID %d", receivedPid)

	if err := t.checkBackgroundCommandExecutablePath(receivedPid); err != nil {
		return err
	}

	logger.Successf("✓ Expected executable path found for process with PID %d", receivedPid)

	if t.SuccessMessage != "" {
		logger.Successf("%s", t.SuccessMessage)
	}

	return nil
}

func (t *BackgroundCommandResponseTestCase) checkBackgroundCommandExecutablePath(pid int) error {
	bgProcess, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("Failed to extract process information on process with PID %d: %s", pid, err)
	}

	receivedExecutablePath, err := bgProcess.Exe()
	if err != nil {
		return fmt.Errorf("Failed to extract arguments for process with PID %d: %s", pid, err)
	}

	cmdlineArgs, err := shlex.Split(t.Command)
	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Failed to extract executable path for command %s: %s", t.Command, err))
	}

	argv0 := cmdlineArgs[0]
	expectedExecutableAbsPath, expectedExecutableResolvedSymlinkPath :=
		utils.MustGetExecutablePathAndResolvedSymlinkForCommand(argv0)

	// If the received executable path is neither the absolute path or the resolved symlink, raise an error
	if !slices.Contains(
		[]string{expectedExecutableAbsPath, expectedExecutableResolvedSymlinkPath},
		receivedExecutablePath,
	) {
		// The error message should contain the unresolved symlink
		return fmt.Errorf("Expected executable path for %s to be %q, got %q", argv0, expectedExecutableAbsPath, receivedExecutablePath)
	}

	return nil
}
