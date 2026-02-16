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

func testFA6(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)
	workingDirPath, err := CreateRandomDirInTmp(stageHarness)
	if err != nil {
		return err
	}
	shell.SetWorkingDirectory(workingDirPath)

	// Create files with nested common prefixes for partial completions
	// e.g., data_abc, data_abc_xyz, data_abc_xyz_123
	basePrefix := fmt.Sprintf("data_%d_", random.RandomInt(1, 100))
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	fileNames := []string{}
	currentPrefix := basePrefix
	for _, word := range randomWords {
		fileName := currentPrefix + word
		currentPrefix = fileName + "_"
		fileNames = append(fileNames, fileName)
		filePath := filepath.Join(workingDirPath, fileName)
		if err := WriteFileWithTeardown(stageHarness, filePath, "", 0644); err != nil {
			return err
		}
	}

	logger.UpdateLastSecondaryPrefix("setup")
	logger.Infof("Available files:")
	for _, fileName := range fileNames {
		logger.Infof("- %s", fileName)
	}
	logger.ResetSecondaryPrefixes()

	if err := asserter.StartShellAndAssertPrompt(false); err != nil {
		return err
	}

	command := GetRandomCommandSuitableForFile()

	// Build expected reflections with command prefix
	expectedReflections := make([]string, len(fileNames))
	for i, fileName := range fileNames {
		expectedReflections[i] = fmt.Sprintf("%s %s", command, fileName)
	}

	// Build raw inputs: first is "command basePrefix", rest are just "_"
	rawInputs := make([]string, len(fileNames))
	rawInputs[0] = fmt.Sprintf("%s %s", command, basePrefix)
	for i := 1; i < len(rawInputs); i++ {
		rawInputs[i] = "_"
	}

	err = test_cases.PartialCompletionsTestCase{
		RawInputs:           rawInputs,
		ExpectedReflections: expectedReflections,
		SuccessMessage:      fmt.Sprintf("Received all partial completions for %q", fileNames[len(fileNames)-1]),
		SkipPromptAssertion: true,
	}.Run(asserter, shell, logger)

	if err != nil {
		return err
	}

	return logAndQuit(asserter, nil)
}
