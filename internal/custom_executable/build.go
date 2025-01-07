package custom_executable

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"
)

func ReplaceAndBuild(outputPath, randomString string) error {
	// Our executable contains a placeholder for the random string
	// The placeholder is <<<RANDOM>>>
	// We will replace the placeholder with the random string
	// The random string HAS to be the same length as the placeholder
	if len(randomString) != 10 {
		return fmt.Errorf("CodeCrafters Internal Error: randomString length must be 10")
	}

	// Copy the custom_executable to the output path
	command := fmt.Sprintf("cp %s %s", path.Join(os.Getenv("TESTER_DIR"), "custom_executable"), outputPath)
	copyCmd := exec.Command("sh", "-c", command)
	copyCmd.Stdout = io.Discard
	copyCmd.Stderr = io.Discard
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: cp failed: %w", err)
	}

	// Replace the placeholder with the random string
	// We can run the executable now, it will work as expected
	err := addSecretCodeToExecutable(outputPath, randomString)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: adding secret code to executable failed: %w", err)
	}

	return nil
}

func CopyExecutable(sourcePath, destinationPath string, logger *logger.Logger) error {
	// Copy the source executable to the destination path
	logger.Infof("Copying %s to %s", sourcePath, destinationPath)
	command := fmt.Sprintf("cp %s %s", sourcePath, destinationPath)
	copyCmd := exec.Command("sh", "-c", command)
	copyCmd.Stdout = io.Discard
	copyCmd.Stderr = io.Discard
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: cp failed: %w", err)
	}
	return nil
}

func CopyExecutableToMultiplePaths(sourcePath string, destinationPaths []string, logger *logger.Logger) error {
	for _, destinationPath := range destinationPaths {
		if err := CopyExecutable(sourcePath, destinationPath, logger); err != nil {
			return err
		}
	}
	return nil
}

func CreateExecutable(randomString, outputPath string) error {
	// Call the replaceAndBuild function
	return ReplaceAndBuild(outputPath, randomString)
}

func addSecretCodeToExecutable(filePath, randomString string) error {
	LENGTH := 10
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: read file failed: %w", err)
	}
	placeholderIndex := strings.Index(string(data), "<<RANDOM>>")
	if placeholderIndex == -1 {
		return fmt.Errorf("CodeCrafters Internal Error: <<RANDOM>> not found in file")
	}
	bytes := copy(data[placeholderIndex:placeholderIndex+LENGTH], []byte(randomString))
	if bytes != LENGTH {
		return fmt.Errorf("CodeCrafters Internal Error: copy failed")
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: write file failed: %w", err)
	}
	return nil
}
