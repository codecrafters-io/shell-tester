package internal

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPA5(stageHarness *test_case_harness.TestCaseHarness) error {
	stageLogger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := os.MkdirAll("xyz_foo/bar", 0755); err != nil {
		return err
	}
	defer os.RemoveAll("xyz_foo")

	if err := os.WriteFile("xyz_foo/bar/new.txt", []byte{}, 0644); err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := random.RandomElementFromArray([]string{"ls", "stat", "file", "du"})

	initialTypedPrefix := fmt.Sprintf("%s %s", command, "xyz_f")
	reflections := []string{
		fmt.Sprintf("%s xyz_foo/", command),
		fmt.Sprintf("%s xyz_foo/bar/", command),
		fmt.Sprintf("%s xyz_foo/bar/new.txt", command),
	}

	err := test_cases.CommandPartialCompletionsTestCase{
		Inputs:              []string{initialTypedPrefix, "", ""},
		ExpectedReflections: reflections,
		SuccessMessage:      fmt.Sprintf("Received all partial completions for %q", initialTypedPrefix),
		SkipPromptAssertion: true,
	}.Run(asserter, shell, stageLogger)
	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
