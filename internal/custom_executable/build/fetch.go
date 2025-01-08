package custom_executable

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/codecrafters-io/tester-utils/logger"
)

func copyFile(sourcePath, destinationPath string, logger *logger.Logger) error {
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

func CopyFileToMultiplePaths(sourcePath string, destinationPaths []string, logger *logger.Logger) error {
	for _, destinationPath := range destinationPaths {
		if err := copyFile(sourcePath, destinationPath, logger); err != nil {
			return err
		}
	}
	return nil
}

func FetchCustomExecutable(fileName string) string {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return fmt.Sprintf("%s_darwin_arm64", fileName)
		case "amd64":
			return fmt.Sprintf("%s_darwin_amd64", fileName)
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			return fmt.Sprintf("%s_linux_arm64", fileName)
		case "amd64":
			return fmt.Sprintf("%s_linux_amd64", fileName)
		}
	}
	panic(fmt.Sprintf("CodeCrafters Internal Error: Unsupported OS:ARCH: %s:%s", runtime.GOOS, runtime.GOARCH))
}

func CopyExecutable(executableName, outputPath string) error {
	// Copy the custom_executable to the output path
	command := fmt.Sprintf("cp %s %s", path.Join(os.Getenv("TESTER_DIR"), executableName), outputPath)
	copyCmd := exec.Command("sh", "-c", command)
	copyCmd.Stdout = io.Discard
	copyCmd.Stderr = io.Discard
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: cp failed: %w", err)
	}

	return nil
}
