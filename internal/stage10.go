package internal

import (
	"fmt"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCd1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	directory, err := GetRandomDirectory(stageHarness)
	if err != nil {
		return err
	}

	testCase := test_cases.CDAndPWDTestCase{Directory: directory, Response: directory}
	err = testCase.Run(asserter, shell, logger)
	if err != nil {
		return err
	}

	directory = "/non-existing-directory"
	command := fmt.Sprintf("cd %s", directory)

	failureTestCase := test_cases.CommandResponseTestCase{
		Command:        command,
		ExpectedOutput: fmt.Sprintf(`cd: %s: No such file or directory`, directory),
		FallbackPatterns: []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf(`^(can't cd to %s|((bash: )?cd: )?%s: No such file or directory)$`, directory, directory)),
			regexp.MustCompile(`^ash: cd: can't cd to ` + directory + `: No such file or directory$`),
			regexp.MustCompile(`^cd: no such file or directory: ` + directory + `$`),
			regexp.MustCompile(`^dash: [0-9]+: cd: can't cd to ` + directory + `$`),
		},
		SuccessMessage: "✓ Received error message",
	}

	if err := failureTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
