package internal

import (
	"fmt"
	"path/filepath"

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

	commandName := "git"
	// The file need not exist in this stage
	completerPath := filepath.Join(
		"/tmp",
		fmt.Sprintf("%s.py", random.RandomWord()),
	)

	registerTestCase := test_cases.CommandWithNoResponseTestCase{
		// Insert extra spaces in between to prevent byte-copying
		Command:        fmt.Sprintf("complete  -C  '%s'  %s", completerPath, commandName),
		SuccessMessage: "✓ No output found",
	}

	if err := registerTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	listCompletionTestCase := test_cases.CommandResponseTestCase{
		Command:        "complete",
		ExpectedOutput: fmt.Sprintf("complete -C '%s' %s", completerPath, commandName),
		SuccessMessage: "✓ Registered completion found in normalized form",
	}

	if err := listCompletionTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
