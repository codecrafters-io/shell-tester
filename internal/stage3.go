package internal

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testREPL(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	numberOfCommands := random.RandomInt(3, 6)

	if err := startShellAndAssertPrompt(asserter, shell); err != nil {
		return err
	}

	for i := 0; i < numberOfCommands; i++ {
		command := "invalid_command_" + strconv.Itoa(i+1)

		testCase := test_cases.CommandResponseTestCase{
			Command:        command,
			ExpectedOutput: fmt.Sprintf("%s: command not found", command),
			FallbackPatterns: []*regexp.Regexp{
				regexp.MustCompile(fmt.Sprintf(`^bash: %s: command not found$`, command)),
				regexp.MustCompile(fmt.Sprintf(`^%s: command not found$`, command)),
			},
			SuccessMessage: "✓ Received command not found message",
		}

		if err := testCase.Run(asserter, shell, logger); err != nil {
			return err
		}
	}

	// AssertWithPrompt() already makes sure that the prompt is present in the last row
	return logAndQuit(asserter, nil)
}
