package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	commandName := random.RandomElementFromArray([]string{"git", "docker", "systemctl"})

	printCompletionSpecTestCase := test_cases.CommandResponseTestCase{
		Command:        fmt.Sprintf("complete -p %s", commandName),
		ExpectedOutput: fmt.Sprintf("complete: %s: no completion specification", commandName),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^bash: complete: %s: no completion specification$`, regexp.QuoteMeta(commandName))),
		},
		SuccessMessage: "✓ Found missing completion specification",
	}

	if err := printCompletionSpecTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
