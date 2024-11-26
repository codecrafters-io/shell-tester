package internal

import (
	"fmt"
	"os"
	"path"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
)

func assertShellIsRunning(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := test_cases.NewSilentPromptTestCase("$ ")

	if err := testCase.Run(shell, logger); err != nil {
		return fmt.Errorf("Expected shell to print prompt after last command, but it didn't: %v", err)
	}
	return nil
}

// GetRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>/<random-word>/<random-word>`
func GetRandomDirectory() (string, error) {
	randomDir := path.Join("/tmp", random.RandomWord(), random.RandomWord(), random.RandomWord())
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
	}
	return randomDir, nil
}

// GetShortRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
func GetShortRandomDirectory() (string, error) {
	randomDir := path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
	}
	return randomDir, nil
}

func GetRandomString() string {
	// We will use a random numeric string of length = 6
	var result string
	for i := 0; i < 5; i++ {
		result += fmt.Sprintf("%d", random.RandomInt(10, 99))
	}

	return result
}

func GetRandomName() string {
	names := []string{"Alice", "David", "Emily", "James", "Maria"}
	return names[random.RandomInt(0, len(names))]
}
