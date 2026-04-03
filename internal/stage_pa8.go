package internal

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA8(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	randomDir, err := CreateShortRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}

	const command = "git"
	completerPath := path.Join(randomDir, "multiCandidateEnvCompleter")
	if err := custom_executable.CreateMultiCandidateEnvCompleter(
		completerPath,
		"COMP_LINE", "COMP_POINT",
		"git", "c", "git",
		"git c", "5",
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

	// Prefix "c" matches commit, checkout; completer prints them sorted (checkout, commit).
	listed := []string{"checkout", "commit"}
	completions := strings.Join(listed, "  ")
	escaped := make([]string, len(listed))
	for i, w := range listed {
		escaped[i] = regexp.QuoteMeta(w)
	}

	multiCase := test_cases.MultipleCompletionsTestCase{
		RawInput:                      "git c",
		TabCount:                      2,
		CheckForBell:                  true,
		ExpectedCompletionOptionsLine: completions,
		ExpectedCompletionOptionsLineFallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile("^" + strings.Join(escaped, `\s+`) + "$"),
		},
		SkipPromptAssertion: true,
		SuccessMessage:      "✓ Multiple programmable completion options listed",
	}
	if err := multiCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
