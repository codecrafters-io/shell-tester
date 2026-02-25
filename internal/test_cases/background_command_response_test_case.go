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
	asserter.AddAssertion(assertions.SingleLineRegexAssertion{
		ExpectedRegexPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\[\d+\]\s+\d+`),
		},
	})

	if err := asserter.AssertWithPrompt(); err != nil {
		return err
	}

	// We match the values later to produce a verbose error message
	outputLine := asserter.Shell.GetScreenState().GetRow(asserter.GetLastLoggedRowIndex())
	outputText := outputLine.String()

	jobNumberRegexp := regexp.MustCompile(`\[(\d+)\]\s+\d+`)
	matches := jobNumberRegexp.FindStringSubmatch(outputText)

	if len(matches) != 2 {
		// This is because the regex is already matched against in the assertion above, we're just re-running this with capture group
		panic(fmt.Sprintf("Codecrafters Internal Error - Shouldn't be here: Could not parse background launch output: %q", outputText))
	}

	actualJobNumber := matches[1]

	if actualJobNumber != fmt.Sprintf("%d", t.ExpectedJobNumber) {
		return fmt.Errorf("Expected job number to be %d, got %s", t.ExpectedJobNumber, actualJobNumber)
	}

	logger.Successf("%s", t.SuccessMessage)
	return nil
}
