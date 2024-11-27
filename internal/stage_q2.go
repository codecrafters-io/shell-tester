package internal

import (
	"fmt"
	"os"
	"path"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomDir, err := GetShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(randomDir)

	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf("f %d", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("f   %d", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("f's%d", random.RandomInt(1, 100))),
	}

	if err := shell.Start(); err != nil {
		return err
	}

	S, L := getRandomWordsSmallAndLarge(5, 5)
	inputs := []string{
		fmt.Sprintf(`echo "%s %s"`, S[1], L[1]),
		fmt.Sprintf(`echo "%s  %s"  "%s"`, S[2], L[2], S[3]),
		fmt.Sprintf(`echo "%s"  "%s's"  "%s"`, S[3], L[4], S[1]),
		fmt.Sprintf(`cat "%s" "%s" "%s"`, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		fmt.Sprintf("%s %s", S[1], L[1]),
		fmt.Sprintf("%s  %s %s", S[2], L[2], S[3]),
		fmt.Sprintf(`%s %s's %s`, S[3], L[4], S[1]),
		`'single'` + `"double" "double's   single"` + `'single' "double" 'single'`,
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents[:3] {
		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:        testCaseContent.Input,
			ExpectedOutput: testCaseContent.ExpectedOutput,
			SuccessMessage: "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	if err := writeFiles(filePaths, []string{`'single'`, `"double" "double's   single"`, `'single' "double" 'single'` + "\n"}, logger); err != nil {
		return err
	}
	testCase := test_cases.SingleLineExactMatchTestCase{
		Command:        testCaseContents[3].Input,
		ExpectedOutput: testCaseContents[3].ExpectedOutput,
		SuccessMessage: "Received expected response",
	}
	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	return assertShellIsRunning(shell, logger)
}
