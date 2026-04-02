package internal

import (
	"fmt"
	"path/filepath"
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

	commandName := "git"
	completerFileBasename := fmt.Sprintf("%s.py", random.RandomWord())
	completerPath := filepath.Join("/tmp", completerFileBasename)

	// Insert extra spaces in between to prevent byte-copying
	registerCmd := fmt.Sprintf("complete  -C  '%s'  %s", completerPath, commandName)

	registerTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        registerCmd,
		SuccessMessage: "✓ No output found",
	}

	if err := registerTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	listTestCase := test_cases.CommandResponseTestCase{
		Command:        "complete",
		ExpectedOutput: fmt.Sprintf("complete -C '%s' %s", completerPath, commandName),
		// MacOS has old version of bash (3.2) where quotes aren't printed
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					"complete -C %s %s",
					regexp.QuoteMeta(completerPath),
					regexp.QuoteMeta(commandName),
				),
			),
		},
		SuccessMessage: "✓ Registered completion found in normalized form",
	}

	if err := listTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
