package internal

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomDir, err := GetRandomDirectory()
	if err != nil {
		return err
	}

	// Add randomDir to PATH (That is where the my_exe file is created)
	currentPath := os.Getenv("PATH")
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, currentPath))

	if err := shell.Start(); err != nil {
		return err
	}

	inputs := []string{
		`echo 'new line'`,
		`echo new     line`,
		`echo 'new     line'`,
	}
	expectedOutputs := []string{
		"new line",
		"new line",
		"new     line",
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents {
		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:                    testCaseContent.Input,
			ExpectedPattern:            fmt.Sprintf("^%s$", testCaseContent.ExpectedOutput),
			ExpectedPatternExplanation: testCaseContent.ExpectedOutput,
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}

type testCaseContent struct {
	Input          string
	ExpectedOutput string
}

func newTestCaseContent(input string, expectedOutput string) testCaseContent {
	return testCaseContent{
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}

func newTestCaseContents(inputs []string, expectedOutputs []string) []testCaseContent {
	testCases := []testCaseContent{}
	for i, input := range inputs {
		testCases = append(testCases, newTestCaseContent(input, expectedOutputs[i]))
	}
	return testCases
}
