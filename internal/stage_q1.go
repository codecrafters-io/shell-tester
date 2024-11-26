package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
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

	// ToDo: Create dir
	// Randomize dir name, use small words
	// Cleanup files after test
	fileDir := "/tmp/q1"
	writeFiles([]string{
		path.Join(fileDir, "f1"),
		path.Join(fileDir, "f2"),
		path.Join(fileDir, "f3"),
	}, []string{"new line", "new line", "new     line\n"}, logger)

	inputs := []string{
		`echo 'new line'`,
		`echo new     line`,
		`echo 'new     line'`,
		fmt.Sprintf(`cat %s/f1 %s/f2 %s/f3`, fileDir, fileDir, fileDir),
	}
	expectedOutputs := []string{
		"new line",
		"new line",
		"new     line",
		`new line` + `new line` + `new     line`,
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents {
		testCase := test_cases.SingleLineStringMatchTestCase{
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

func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func writeFiles(paths []string, contents []string, logger *logger.Logger) error {
	for i, content := range contents {
		logger.Infof("Writing file %s with content \"%s\"", paths[i], strings.TrimRight(content, "\n"))
		if err := writeFile(paths[i], content); err != nil {
			logger.Errorf("Error writing file %s: %v", paths[i], err)
			return err
		}
	}
	return nil
}

var SMALL_WORDS = []string{"foo", "bar", "baz", "qux", "quz"}
var LARGE_WORDS = []string{"hello", "world", "test", "example", "shell", "script"}

func getRandomWordsSmallAndLarge(smallCount int, largeCount int) ([]string, []string) {
	smallWords := random.RandomElementsFromArray(SMALL_WORDS, smallCount)
	largeWords := random.RandomElementsFromArray(LARGE_WORDS, largeCount)
	return smallWords, largeWords
}

var ADJECTIVES = []string{"cute", "soft", "furry", "tiny", "cozy", "sweet", "warm", "calm"}
