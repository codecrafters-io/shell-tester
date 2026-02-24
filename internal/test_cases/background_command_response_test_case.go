package test_cases

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// BackgroundCommandResponseTestCase launches the given command with an & symbol
// Launching it to the background
// It asserts that the line that follows will match the expected fallback patterns
// It will record the PID of the background job launched
// It asserts the next prompt immediately
type BackgroundCommandResponseTestCase struct {
	Command           string
	launchedJobNumber *int
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

	asserter.AddAssertion(assertions.SingleLineAssertion{
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\[\d+\]\s+\d+`),
		},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	outputLine := asserter.Shell.GetScreenState().GetRow(asserter.GetLastLoggedRowIndex())
	outputText := outputLine.String()

	// Keeping the capture group for PGID as well: we might need it later
	jobNumberRegexp := regexp.MustCompile(`\[(\d+)\]\s+(\d+)`)
	matches := jobNumberRegexp.FindStringSubmatch(outputText)

	if len(matches) != 3 {
		panic(fmt.Sprintf("Codecrafters Internal Error - Shouldn't be here: Could not parse output: %q", outputText))
	}

	jobNumberStr := matches[1]

	jobNumber, err := strconv.Atoi(jobNumberStr)
	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Shouldn't be here: Could not parse job number from output: %q", outputText))
	}

	t.launchedJobNumber = &jobNumber

	logger.Successf("%s", t.SuccessMessage)
	return nil
}

func (t *BackgroundCommandResponseTestCase) GetLaunchedJobNumber() int {
	if t.launchedJobNumber == nil {
		panic("Codecrafters Internal Error - GetLastLaunchJobNumber called without successful run of the test case")
	}
	return *t.launchedJobNumber
}
