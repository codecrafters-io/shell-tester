package internal

import (
	"fmt"
	"path"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA9(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	const command = "git"
	completerPath := path.Join(randomDir, "lcpEnvCompleter")
	if err := custom_executable.CreateLCPEnvCompleter(
		completerPath,
		"COMP_LINE", "COMP_POINT",
		command, command,
	); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	registerCmd := fmt.Sprintf("complete -C %s %s", completerPath, command)
	registerTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        registerCmd,
		SuccessMessage: "✓ Registered command-based completion",
	}
	if err := registerTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	// LCP of checkout + cherry-pick for prefix "c" is "che".
	if err := (test_cases.AutocompleteTestCase{
		RawInput:            "git c",
		ExpectedCompletion:  "git che",
		SkipPromptAssertion: true,
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	// At LCP "che", still ambiguous → bell; line unchanged.
	if err := (test_cases.AutocompleteTestCase{
		PreviousInputOnLine: "git che",
		RawInput:            "",
		ExpectedCompletion:  "git che",
		CheckForBell:        true,
		SkipPromptAssertion: true,
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	// "chec" matches only checkout → full word + trailing space.
	if err := (test_cases.AutocompleteTestCase{
		PreviousInputOnLine: "git che",
		RawInput:            "c",
		ExpectedCompletion:  "git checkout ",
		SkipPromptAssertion: true,
	}).Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
