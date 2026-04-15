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

func testPEX7(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	unsetVariable := "unset_variable_" + strconv.Itoa(random.RandomInt(1, 100))

	echoTestCase := test_cases.CommandResponseTestCase{
		Command:          fmt.Sprintf("echo $%s", unsetVariable),
		ExpectedOutput:   "",
		FallbackPatterns: []*regexp.Regexp{regexp.MustCompile("^$")},
		SuccessMessage:   "✓ Received expected response",
	}
	if err := echoTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
