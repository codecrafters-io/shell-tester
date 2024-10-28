package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testEcho(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	numberOfCommands := random.RandomInt(2, 4)

	if err := shell.Start(); err != nil {
		return err
	}

	for i := 0; i < numberOfCommands; i++ {
		words := strings.Join(random.RandomWords(random.RandomInt(2, 4)), " ")
		command := fmt.Sprintf("echo %s", words)

		testCase := test_cases.SingleLineOutputTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(fmt.Sprintf(`^%s$`, words)),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", words),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}
