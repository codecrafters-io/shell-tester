package internal

import (
	"fmt"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ5(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	_, L := getRandomWordsSmallAndLarge(5, 5)

	inputs := []string{
		fmt.Sprintf(`echo "%s'%s'\\n'%s"`, L[0], L[1], L[2]),
		fmt.Sprintf(`echo "%s\"insidequotes"%s\"`, L[0], L[1]),
		fmt.Sprintf(`echo "%s\"inner"\\"word\"%s"`, L[2], L[3]),
		fmt.Sprintf(`echo "mixed\"quote'%s'\\"`, L[4]),
	}
	expectedOutputs := []string{
		fmt.Sprintf(`%s'%s'\n'%s`, L[0], L[1], L[2]),
		fmt.Sprintf(`%s"insidequotes%s"`, L[0], L[1]),
		fmt.Sprintf(`%s"inner\word"%s`, L[2], L[3]),
		fmt.Sprintf(`mixed"quote'%s'\`, L[4]),
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents {
		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:        testCaseContent.Input,
			ExpectedOutput: testCaseContent.ExpectedOutput,
			SuccessMessage: "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}
