package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/logged_shell_asserter"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

var SMALL_WORDS = []string{"foo", "bar", "baz", "qux", "quz"}
var LARGE_WORDS = []string{"hello", "world", "test", "example", "shell", "script"}

const CUSTOM_LS_COMMAND = "ls"
const CUSTOM_CAT_COMMAND = "cat"

// getRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>/<random-word>/<random-word>`
func getRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	randomDir := path.Join("/tmp", random.RandomWord(), random.RandomWord(), random.RandomWord())
	if err := CreateDirectory(randomDir, 0755); err != nil {
		return "", err
	}

	// Automatically cleanup the directory when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		grandParentDir := path.Dir(path.Dir(randomDir))
		cleanupDirectories([]string{grandParentDir})
	})

	return randomDir, nil
}

func CreateDirectory(path string, fileMode os.FileMode) error {
	if err := os.MkdirAll(path, fileMode); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: Error creating directory %s: %v", path, err)
	}

	return nil
}

func getRandomInvalidCommand() string {
	return "invalid_" + random.RandomWord() + "_command"
}

func getRandomInvalidCommands(n int) []string {
	words := random.RandomWords(n)
	invalidCommands := make([]string, n)

	for i := 0; i < n; i++ {
		invalidCommands[i] = "invalid_" + words[i] + "_command"
	}

	return invalidCommands
}

// getShortRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
func getShortRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	randomDir := path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
	if err := CreateDirectory(randomDir, 0755); err != nil {
		return "", err
	}

	// Automatically cleanup the directory when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		cleanupDirectories([]string{randomDir})
	})

	return randomDir, nil
}

// TODO: Refactor this to use getShortRandomDirectory internally
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

func getRandomString() string {
	// We will use a random numeric string of length = 6
	var result string
	for i := 0; i < 5; i++ {
		result += fmt.Sprintf("%d", random.RandomInt(10, 99))
	}

	return result
}

func getRandomName() string {
	names := []string{"Alice", "David", "Emily", "James", "Maria"}
	return names[random.RandomInt(0, len(names))]
}

// WriteFile writes a file to the given path with the given content
func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func ChangeFilePermissions(path string, fileMode os.FileMode) error {
	return os.Chmod(path, fileMode)
}

// writeFiles writes a list of files to the given paths with the given contents
func writeFiles(paths []string, contents []string, logger *logger.Logger) error {
	for i, content := range contents {
		logger.Infof("Writing file \"%s\" with content \"%s\"", paths[i], strings.TrimRight(content, "\n"))
		if err := WriteFile(paths[i], content); err != nil {
			logger.Errorf("Error writing file %s: %v", paths[i], err)
			return err
		}
	}
	return nil
}

func logAndQuit(asserter *logged_shell_asserter.LoggedShellAsserter, err error) error {
	asserter.LogRemainingOutput()
	return err
}

func SetUpCustomCommands(stageHarness *test_case_harness.TestCaseHarness, shell *shell_executable.ShellExecutable, commands []string) (string, error) {
	executableDir, err := getRandomDirectory(stageHarness)
	if err != nil {
		return "", err
	}
	// Add the random directory to PATH
	// (where the custom executable is copied to)
	shell.AddToPath(executableDir)

	for _, command := range commands {
		switch command {
		case "ls":
			customLsPath := path.Join(executableDir, CUSTOM_LS_COMMAND)
			err = custom_executable.CreateLsExecutable(customLsPath)
			if err != nil {
				return "", err
			}
		case "cat":
			customCatPath := path.Join(executableDir, CUSTOM_CAT_COMMAND)
			err = custom_executable.CreateCatExecutable(customCatPath)
			if err != nil {
				return "", err
			}
		}
	}

	return executableDir, nil
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}
