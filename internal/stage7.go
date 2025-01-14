package internal

import (
	"fmt"
	"os"
	"path/filepath"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	// Add the random directory to PATH (where the my_exe file is created)
	randomDir, err := getRandomDirectory(stageHarness)
	if err != nil {
		return err
	}
	randomWords := random.RandomElementsFromArray(SMALL_WORDS, 3)
	childDir1 := filepath.Join(randomDir, randomWords[0])
	childDir2 := filepath.Join(randomDir, randomWords[1])

	path := os.Getenv("PATH")
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	shell.Setenv("PATH", fmt.Sprintf("%s:%s:%s:%s", childDir1, childDir2, randomDir, path))
	asserter := logged_shell_asserter.NewLoggedShellAsserter(shell)

	customExecutablePath := filepath.Join(randomDir, "my_exe")
	err = custom_executable.CreateSignaturePrinterExecutable(getRandomString(), customExecutablePath)
	if err != nil {
		return err
	}

	if err := asserter.StartShellAndAssertPrompt(); err != nil {
		return err
	}

	availableExecutables := []string{"cat", "cp", "mkdir", "my_exe"}

	for _, executable := range availableExecutables {
		testCase := test_cases.TypeOfCommandTestCase{
			Command: executable,
		}

		var expectedPath = ""
		if executable == "my_exe" {
			expectedPath = customExecutablePath
		}

		if err := testCase.RunForExecutable(asserter, shell, logger, expectedPath); err != nil {
			return err
		}
	}

	invalidCommands := getRandomInvalidCommands(2)

	for _, invalidCommand := range invalidCommands {
		testCase := test_cases.TypeOfCommandTestCase{
			Command: invalidCommand,
		}
		if err := testCase.RunForInvalidCommand(asserter, shell, logger); err != nil {
			return err
		}
	}

	logger.UpdateSecondaryPrefix("Setup")
	// randomDir is on PATH
	logger.Infof("mkdir -p %s", childDir1)
	CreateDirectory(childDir1, 0755)
	logger.Infof("mkdir -p %s", childDir2)
	CreateDirectory(childDir2, 0755)
	file1 := filepath.Join(childDir1, randomWords[2])
	file2 := filepath.Join(childDir2, randomWords[2])
	logger.Infof("touch %s", file1)
	if err := WriteFile(file1, "hello"); err != nil {
		return err
	}
	logger.Infof("touch %s", file2)
	if err := WriteFile(file2, "world"); err != nil {
		return err
	}

	shell.SendCommand("echo $PATH")
	shell.SendCommand(fmt.Sprintf("stat %s", file1))
	shell.SendCommand(fmt.Sprintf("stat %s", file2))

	testCase := test_cases.TypeOfCommandTestCase{
		Command: randomWords[2],
	}
	if err := testCase.RunForInvalidCommand(asserter, shell, logger); err != nil {
		return err
	}
	logger.Infof("chmod 755 %s", file1)
	if err := ChangeFilePermissions(file1, 0755); err != nil {
		return err
	}

	if err := testCase.RunForExecutable(asserter, shell, logger, file1); err != nil {
		return err
	}

	logger.ResetSecondaryPrefix()
	return logAndQuit(asserter, nil)
}
