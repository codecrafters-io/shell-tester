package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ3(stageHarness *test_case_harness.TestCaseHarness) error {
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

	fileDir := "/tmp/"
	fileDir = filepath.Join(fileDir, random.RandomElementFromArray([]string{"foo", "bar", "baz"}))
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, 0755)
	}

	writeFiles([]string{
		path.Join(fileDir, "f1"),
		path.Join(fileDir, "f2"),
		path.Join(fileDir, "f3"),
	}, []string{`Hello\\nWorld`, `\\\\`, `'single' '\n\n\n\''` + "\n"}, logger)

	_, L := getRandomWordsSmallAndLarge(5, 5)
	inputs := []string{
		`echo "before\   after"`,
		fmt.Sprintf(`echo %s\ \ \ \ \ \ %s`, L[0], L[1]),
		fmt.Sprintf(`echo %s\n%s`, L[2], L[3]),
		fmt.Sprintf(`cat %s/f1 %s/f2 %s/f3`, fileDir, fileDir, fileDir),
	}
	expectedOutputs := []string{
		`before\   after`,
		fmt.Sprintf("%s      %s", L[0], L[1]),
		fmt.Sprintf("%sn%s", L[2], L[3]),
		`Hello\\nWorld` + `\\\\` + `'single' '\n\n\n\''`,
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
