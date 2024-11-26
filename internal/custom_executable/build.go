package custom_executable

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
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
	command = fmt.Sprintf("echo -n \"%s\" | dd of=%s bs=1 seek=$((0x95070 + 4)) conv=notrunc", randomString, outputPath)
	buildCmd := exec.Command("sh", "-c", command)
	buildCmd.Stdout = io.Discard
	buildCmd.Stderr = io.Discard
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: dd replace failed: %w", err)
	}

	return nil
}

func CopyExecutable(sourcePath, destinationPath string) error {
	// Copy the source executable to the destination path
	command := fmt.Sprintf("cp %s %s", sourcePath, destinationPath)
	copyCmd := exec.Command("sh", "-c", command)
	copyCmd.Stdout = io.Discard
	copyCmd.Stderr = io.Discard
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: cp failed: %w", err)
	}
	return nil
}

func CopyExecutableToMultiplePaths(sourcePath string, destinationPaths []string) error {
	for _, destinationPath := range destinationPaths {
		if err := CopyExecutable(sourcePath, destinationPath); err != nil {
			return err
		}
	}
	return nil
}

func CreateExecutable(randomString, outputPath string) error {
	// Call the replaceAndBuild function
	return ReplaceAndBuild(outputPath, randomString)
}
