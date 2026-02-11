package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// getRandomDirectory creates a random directory in /tmp,
// creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>/<random-word>/<random-word>`
// If performCleanup is true, the directory will be cleaned up
// when the test is completed
// The total possible directories is 10^3 = 1000
// This can be used without cleanup in most cases
func getRandomDirectory(stageHarness *test_case_harness.TestCaseHarness, performCleanup bool) (string, error) {
	randomDir := path.Join("/tmp", random.RandomWord(), random.RandomWord(), random.RandomWord())
	for {
		if _, err := os.Stat(randomDir); os.IsNotExist(err) {
			if err := os.MkdirAll(randomDir, 0755); err != nil {
				return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
			}
			break
		}
		randomDir = path.Join("/tmp", random.RandomWord(), random.RandomWord(), random.RandomWord())
	}

	// Automatically cleanup the directory when the test is completed, if requested
	if performCleanup {
		stageHarness.RegisterTeardownFunc(func() {
			grandParentDir := path.Dir(path.Dir(randomDir))
			cleanupDirectories([]string{grandParentDir})
		})
	}

	return randomDir, nil
}

func GetRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	return getRandomDirectory(stageHarness, true)
}

// GetShortRandomDirectory creates a random directory in /tmp,
// creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
// Cleanup is performed automatically, and as the total possible directories
// is very small, this should not be used without cleanup
func GetShortRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
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

func GetShortRandomDirectories(stageHarness *test_case_harness.TestCaseHarness, n int) ([]string, error) {
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

// CreateRandomFileInDir creates a random file inside the given directory
// If extension is non empty, it is used as the filename extension
// Returns file basename, contents and error encountered (if any) during creation
func CreateRandomFileInDir(stageHarness *test_case_harness.TestCaseHarness, dirPath string, extension string, filemode os.FileMode) (string, string, error) {
	fileBaseName := fmt.Sprintf("%s-%d", random.RandomWord(), random.RandomInt(1, 100))
	if extension != "" {
		fileBaseName += "." + extension
	}
	filePath := filepath.Join(dirPath, fileBaseName)
	contents := random.RandomString()

	if err := os.WriteFile(filePath, []byte(contents), filemode); err != nil {
		return "", "", err
	}

	stageHarness.RegisterTeardownFunc(func() {
		os.Remove(filePath)
	})

	return fileBaseName, contents, nil
}

// MkdirAllWithTeardown is a wrapper over os.Mkdir that registers teardown to delete the directory using the harness
// It will delete every directory in the hierarchy that was created by it
func MkdirAllWithTeardown(
	stageHarness *test_case_harness.TestCaseHarness,
	dirPath string,
	permissions os.FileMode,
) error {
	abs, err := filepath.Abs(dirPath)
	if err != nil {
		return err
	}

	toRemove := abs
	for {
		parent := filepath.Dir(toRemove)
		if parent == toRemove {
			break
		}
		if _, err := os.Stat(parent); os.IsNotExist(err) {
			toRemove = parent
		} else {
			break
		}
	}

	if err := os.MkdirAll(abs, permissions); err != nil {
		return err
	}

	stageHarness.RegisterTeardownFunc(func() {
		_ = os.RemoveAll(toRemove)
	})

	return nil
}

func WriteFileWithTeardown(stageHarness *test_case_harness.TestCaseHarness, filePath string, contents string, permissions os.FileMode) error {
	if err := os.WriteFile(filePath, []byte(contents), permissions); err != nil {
		return err
	}

	stageHarness.RegisterTeardownFunc(func() {
		os.Remove(filePath)
	})

	return nil
}

type WriteFileSpec struct {
	FilePath    string
	FileContent string
	Permission  os.FileMode
}

func WriteFilesWithTearDown(stageHarness *test_case_harness.TestCaseHarness, writeFileSpecs []WriteFileSpec) error {
	for _, spec := range writeFileSpecs {
		if err := WriteFileWithTeardown(stageHarness, spec.FilePath, spec.FileContent, spec.Permission); err != nil {
			return err
		}
	}
	return nil
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
