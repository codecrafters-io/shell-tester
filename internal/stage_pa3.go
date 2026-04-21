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

func testPA3(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	if err := asserter.StartShellAndAssertPrompt(true); err != nil {
		return err
	}

	// Scripts need not exist; paths must differ (two random basenames under /tmp).
	words := random.RandomWords(2)
	gitCompleterPath := filepath.Join("/tmp", fmt.Sprintf("%s.py", words[0]))
	dockerCompleterPath := filepath.Join("/tmp", fmt.Sprintf("%s.py", words[1]))

	registerGitTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        fmt.Sprintf("complete  -C  '%s'  git", gitCompleterPath),
		SuccessMessage: "✓ No output found",
	}
	if err := registerGitTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	registerDockerTestCase := test_cases.CommandWithNoResponseTestCase{
		Command:        fmt.Sprintf("complete  -C  '%s'  docker", dockerCompleterPath),
		SuccessMessage: "✓ No output found",
	}
	if err := registerDockerTestCase.Run(asserter, shell, logger, false); err != nil {
		return err
	}

	printGitSpecTestCase := test_cases.CommandResponseTestCase{
		Command:        "complete -p git",
		ExpectedOutput: fmt.Sprintf("complete -C '%s' git", gitCompleterPath),
		SuccessMessage: "✓ Registered git completion found in normalized form",
	}
	if err := printGitSpecTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	printDockerSpecTestCase := test_cases.CommandResponseTestCase{
		Command:        "complete -p docker",
		ExpectedOutput: fmt.Sprintf("complete -C '%s' docker", dockerCompleterPath),
		SuccessMessage: "✓ Registered docker completion found in normalized form",
	}
	if err := printDockerSpecTestCase.Run(asserter, shell, logger); err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
