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

func testPA10(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	commandName := "git"
	completerFileBasename := fmt.Sprintf("%s.py", random.RandomWord())
	completerPath := filepath.Join("/tmp", completerFileBasename)

	registerCmd := fmt.Sprintf("complete  -C  '%s'  %s", completerPath, commandName)
	if err := (test_cases.CommandWithNoResponseTestCase{
		Command:        registerCmd,
		SuccessMessage: "✓ No output found",
	}).Run(asserter, shell, logger, false); err != nil {
		return err
	}

	unregisterCmd := fmt.Sprintf("complete -r %s", commandName)
	if err := (test_cases.CommandWithNoResponseTestCase{
		Command:        unregisterCmd,
		SuccessMessage: "✓ No output from complete -r",
	}).Run(asserter, shell, logger, false); err != nil {
		return err
	}

	if err := (test_cases.AutocompleteTestCase{
		RawInput:            commandName + " ",
		ExpectedCompletion:  commandName + " ",
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
