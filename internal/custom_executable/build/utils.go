package custom_executable

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func addSecretCodeToExecutable(filePath, randomString string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: read file failed: %w", err)
	}
	placeholder := []byte("<<RANDOM>>")
	if !bytes.Contains(data, placeholder) {
		return fmt.Errorf("CodeCrafters Internal Error: placeholder %q not found in %s", placeholder, filePath)
	}

	newData := bytes.ReplaceAll(data, placeholder, []byte(randomString))
	if err := os.WriteFile(filePath, newData, 0644); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: write file failed: %w", err)
	}
	return nil
}

func reSignExecutableDarwinArm64(executablePath string) error {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		return nil
	}
	if err := exec.Command("codesign", "--remove-signature", executablePath).Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: removing signature from executable failed: %w", err)
	}
	if err := exec.Command("codesign", "-s", "-", executablePath).Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: signing executable failed: %w", err)
	}
	return nil
}

// fetchCustomExecutableForOSAndArch is a helper function to
// fetch the correct custom executable for the current OS and architecture
// It just returns the name of the executable, without the path
// Path is added in the caller, which is expected to be the top level directory
func fetchCustomExecutableForOSAndArch(fileName string) string {
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

func copyExecutable(executableName, outputPath string) error {
	// Copy the custom_executable to the output path
	command := fmt.Sprintf("cp %s %s", path.Join(os.Getenv("TESTER_DIR"), "built_executables", executableName), outputPath)
	copyCmd := exec.Command("sh", "-c", command)
	copyCmd.Stdout = io.Discard
	copyCmd.Stderr = io.Discard
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: cp failed: %w", err)
	}

	return nil
}

// createExecutableForOSAndArch is a helper function to
// fetch the correct custom executable for the current OS and architecture
// and places it in the outputPath
func createExecutableForOSAndArch(executableName, outputPath string) error {
	fileName := fetchCustomExecutableForOSAndArch(executableName)

	// Copy the base executable from archive location to user's executable path
	err := copyExecutable(fileName, outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: copying executable %s failed: %w", fileName, err)
	}

	return nil
}
