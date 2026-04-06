package custom_executable

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

// secretSlotByteLength is the required byte length for the replacement value at slot n (same as len(<<RANDOM_n>>)).
func secretSlotByteLength(n int) int {
	return len(secretPlaceholder(n))
}

// prepareSecretPatchedExecutable copies the prebuilt binary baseExecutableName (under built_executables),
// replaces <<RANDOM_1>>..<<RANDOM_N>> with secrets[0]..secrets[N-1] (each value may be shorter; it is right-padded
// with spaces to SecretSlotByteLen(i)), and re-signs on darwin/arm64.
func prepareSecretPatchedExecutable(baseExecutableName, outputPath string, secrets []string) error {
	if err := createExecutableForOSAndArch(baseExecutableName, outputPath); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: copying executable %s failed: %w", baseExecutableName, err)
	}

	if err := applyNumberedSecretPatches(outputPath, secrets); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: patching executable %s failed: %w", baseExecutableName, err)
	}

	return reSignExecutableDarwinARM64(outputPath)
}

// secretPlaceholder returns the embedded token for slot n (1-based: <<RANDOM_1>>, <<RANDOM_2>>, ...).
// On disk, replacements are padded to this width so the binary size stays fixed (see PadSecretForSlot).
func secretPlaceholder(n int) string {
	if n < 1 {
		panic("secret slot n must be >= 1")
	}
	return fmt.Sprintf("<<RANDOM_%d>>", n)
}

func applyNumberedSecretPatches(filePath string, secrets []string) error {
	if len(secrets) == 0 {
		return fmt.Errorf("CodeCrafters Internal Error: at least one secret patch is required")
	}

	fileContents, readErr := os.ReadFile(filePath)

	if readErr != nil {
		return fmt.Errorf("CodeCrafters Internal Error: read file failed: %w", readErr)
	}

	for secretIndex, secretValue := range secrets {
		slotNumber := secretIndex + 1
		placeholderString := secretPlaceholder(slotNumber)
		placeholderBytes := []byte(placeholderString)

		paddedSecret, padErr := padSecretForSlot(slotNumber, secretValue)
		if padErr != nil {
			return padErr
		}

		if !bytes.Contains(fileContents, placeholderBytes) {
			return fmt.Errorf("CodeCrafters Internal Error: placeholder %q not found in %s", placeholderString, filePath)
		}

		fileContents = bytes.ReplaceAll(fileContents, placeholderBytes, []byte(paddedSecret))
	}

	if writeErr := os.WriteFile(filePath, fileContents, 0644); writeErr != nil {
		return fmt.Errorf("CodeCrafters Internal Error: write file failed: %w", writeErr)
	}

	return nil
}

// padSecretForSlot returns s right-padded with ASCII spaces to SecretSlotByteLen(n), the width embedded in the binary.
// Logical values may be shorter; comparers should trim trailing spaces (e.g. strings.TrimRight(s, " ")) when needed.
func padSecretForSlot(n int, s string) (string, error) {
	max := secretSlotByteLength(n)
	if len(s) > max {
		return "", fmt.Errorf("CodeCrafters Internal Error: value for slot %d exceeds %d bytes", n, max)
	}
	return s + strings.Repeat(" ", max-len(s)), nil
}

func reSignExecutableDarwinARM64(outputPath string) error {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		return nil
	}

	if err := exec.Command("codesign", "--remove-signature", outputPath).Run(); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: removing signature from executable failed: %w", err)
	}

	if err := exec.Command("codesign", "-s", "-", outputPath).Run(); err != nil {
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
