package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// getRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>/<random-word>/<random-word>`
func getRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	randomDir := path.Join("/tmp", random.RandomWord(), random.RandomWord(), random.RandomWord())
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
	}

	// Automatically cleanup the directory when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		grandParentDir := path.Dir(path.Dir(randomDir))
		cleanupDirectories([]string{grandParentDir})
	})

	return randomDir, nil
}

// getShortRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
func getShortRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	randomDir := path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
	}

	// Automatically cleanup the directory when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		cleanupDirectories([]string{randomDir})
	})

	return randomDir, nil
}

func getShortRandomDirectories(stageHarness *test_case_harness.TestCaseHarness, n int) ([]string, error) {
	directoryNames := random.RandomElementsFromArray(SMALL_WORDS, n)
	randomDirs := make([]string, n)
	for i := 0; i < n; i++ {
		randomDir := path.Join("/tmp", directoryNames[i])
		if err := os.MkdirAll(randomDir, 0755); err != nil {
			return nil, fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
		}
		randomDirs[i] = randomDir
	}

	// Automatically cleanup the directories when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		cleanupDirectories(randomDirs)
	})

	return randomDirs, nil
}

// writeFile writes a file to the given path with the given content
func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// writeFiles writes a list of files to the given paths with the given contents
func writeFiles(paths []string, contents []string, logger *logger.Logger) error {
	for i, content := range contents {
		logger.UpdateSecondaryPrefix("Setup")
		logger.Infof("echo -n %q > %q", strings.TrimRight(content, "\n"), paths[i])
		logger.ResetSecondaryPrefix()

		if err := writeFile(paths[i], content); err != nil {
			logger.Errorf("Error writing file %s: %v", paths[i], err)
			return err
		}
	}
	return nil
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}
