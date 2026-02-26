package test_cases

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/assertions"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
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

	// We first match against the format
	asserter.AddAssertion(assertions.SingleLineAssertion{
		ExpectedOutput: fmt.Sprintf("[%d] <PID>", t.ExpectedJobNumber),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^\[%d\] \d+$`, t.ExpectedJobNumber)),
		},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	if t.SuccessMessage != "" {
		logger.Successf("%s", t.SuccessMessage)
	}

	return nil
}
