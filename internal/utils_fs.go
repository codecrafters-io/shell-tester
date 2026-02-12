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

// CreateRandomDirIn creates a random directory in rootDir,
// creates the directories and returns the full path
// directory is of the form `rootDir/<random-word>/<random-word>/<random-word>`
// If performCleanup is true, the directory will be cleaned up
// when the test is completed
// The total possible directories is 10^3 = 1000
// This can be used without cleanup in most cases
func CreateRandomDirIn(stageHarness *test_case_harness.TestCaseHarness, rootDir string) (string, error) {
	randomDir := path.Join(rootDir, random.RandomWord(), random.RandomWord(), random.RandomWord())
	for {
		if _, err := os.Stat(randomDir); os.IsNotExist(err) {
			if err := os.MkdirAll(randomDir, 0755); err != nil {
				return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
			}
			break
		}
		randomDir = path.Join(rootDir, random.RandomWord(), random.RandomWord(), random.RandomWord())
	}

	// Automatically cleanup the directory when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		grandParentDir := path.Dir(path.Dir(randomDir))
		cleanupDirectories([]string{grandParentDir})
	})

	return randomDir, nil
}

func CreateRandomDirInTmp(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	return CreateRandomDirIn(stageHarness, "/tmp")
}

// CreateShortRandomDirInTmp creates a random directory in /tmp,
// creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
// Cleanup is performed automatically, and as the total possible directories
// is very small, this should not be used without cleanup
func CreateShortRandomDirInTmp(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	seen := make(map[string]bool)
	randomDir := path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
	for {
		seen[randomDir] = true
		if _, err := os.Stat(randomDir); os.IsNotExist(err) {
			if err := os.MkdirAll(randomDir, 0755); err != nil {
				return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
			}
			break
		}
		randomDir = path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
		if len(seen) == len(SMALL_WORDS) {
			// We've seen all possible directories, so we return randomDir
			// We are okay with returning an used directory here
			// TODO: Possibly return error here instead ?
			break
		}
	}

	// Automatically cleanup the directory when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		cleanupDirectories([]string{randomDir})
	})

	return randomDir, nil
}

func CreateShortRandomDirsInTmp(stageHarness *test_case_harness.TestCaseHarness, n int) ([]string, error) {
	if n > len(SMALL_WORDS) {
		panic(fmt.Sprintf("CodeCrafters internal error. Number of directories to create is greater than the number of possible directories: %d", n))
	}

	directoryNames := random.RandomElementsFromArray(SMALL_WORDS, n)
	randomDirs := make([]string, n)
	for i := range n {
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
func writeFile(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

func appendFile(filePath string, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	// ensure the file is flushed to disk immediately
	// we don't care about the error here
	file.Sync()

	return nil
}

// writeFiles writes a list of files to the given paths with the given contents
func writeFiles(paths []string, contents []string, logger *logger.Logger) error {
	for i, content := range contents {
		logger.UpdateLastSecondaryPrefix("setup")

		if strings.HasSuffix(content, "\n") {
			if strings.Count(content, "\n") == 1 { // Trailing newline only
				logger.Infof("echo %q > %q", content[:len(content)-1], paths[i])
			} else { // Newline(s) in the middle as well
				logger.Infof("echo -e %q > %q", content[:len(content)-1], paths[i])
			}
		} else {
			if strings.Contains(content, "\n") {
				logger.Infof("echo -ne %q > %q", content, paths[i])
			} else {
				logger.Infof("echo -n %q > %q", content, paths[i])
			}
		}

		logger.ResetSecondaryPrefixes()

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
