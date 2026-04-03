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

type secretPatchedExecutable int

const (
	secretPatchedSignaturePrinter secretPatchedExecutable = iota
	secretPatchedSingleCompleter
	secretPatchedNoStdoutStderrCompleter
	secretPatchedEnvContextCompleter
)

// secretPlaceholder returns the embedded token for slot n (1-based: <<RANDOM_1>>, <<RANDOM_2>>, ...).
// Replacement strings must be exactly len(secretPlaceholder(n)) bytes so the binary size stays fixed.
func secretPlaceholder(n int) string {
	if n < 1 {
		panic("secret slot n must be >= 1")
	}
	return fmt.Sprintf("<<RANDOM_%d>>", n)
}

// SecretSlotByteLen is the required byte length for the replacement value at slot n (same as len(<<RANDOM_n>>)).
func SecretSlotByteLen(n int) int {
	return len(secretPlaceholder(n))
}

func applyNumberedSecretPatches(filePath string, secrets []string) error {
	if len(secrets) == 0 {
		return fmt.Errorf("CodeCrafters Internal Error: at least one secret patch is required")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: read file failed: %w", err)
	}
	for i, secret := range secrets {
		slot := i + 1
		ph := secretPlaceholder(slot)
		phb := []byte(ph)
		if len(secret) != len(ph) {
			return fmt.Errorf("CodeCrafters Internal Error: secret for %q must be %d bytes, got %d", ph, len(ph), len(secret))
		}
		if !bytes.Contains(data, phb) {
			return fmt.Errorf("CodeCrafters Internal Error: placeholder %q not found in %s", ph, filePath)
		}
		data = bytes.ReplaceAll(data, phb, []byte(secret))
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: write file failed: %w", err)
	}
	return nil
}

// prepareSecretPatchedExecutable copies the built binary for the current OS/arch, replaces <<RANDOM_1>>..N
// with secrets[i] for each slot, and re-signs on darwin/arm64.
func prepareSecretPatchedExecutable(kind secretPatchedExecutable, outputPath string, secrets []string) error {
	var baseName string
	switch kind {
	case secretPatchedSignaturePrinter:
		baseName = "signature_printer"
	case secretPatchedSingleCompleter:
		baseName = "single_completer"
	case secretPatchedNoStdoutStderrCompleter:
		baseName = "no_stdout_stderr_completer"
	case secretPatchedEnvContextCompleter:
		baseName = "env_context_completer"
	default:
		return fmt.Errorf("CodeCrafters Internal Error: unknown patched executable kind")
	}

	err := createExecutableForOSAndArch(baseName, outputPath)
	if err != nil {
		switch kind {
		case secretPatchedSignaturePrinter:
			return fmt.Errorf("CodeCrafters Internal Error: copying executable failed: %w", err)
		case secretPatchedSingleCompleter:
			return fmt.Errorf("CodeCrafters Internal Error: copying single_completer failed: %w", err)
		case secretPatchedNoStdoutStderrCompleter:
			return fmt.Errorf("CodeCrafters Internal Error: copying no_stdout_stderr_completer failed: %w", err)
		case secretPatchedEnvContextCompleter:
			return fmt.Errorf("CodeCrafters Internal Error: copying env_context_completer failed: %w", err)
		}
	}

	err = applyNumberedSecretPatches(outputPath, secrets)
	if err != nil {
		switch kind {
		case secretPatchedSignaturePrinter:
			return fmt.Errorf("CodeCrafters Internal Error: adding secret code to executable failed: %w", err)
		case secretPatchedSingleCompleter:
			return fmt.Errorf("CodeCrafters Internal Error: patching single_completer failed: %w", err)
		case secretPatchedNoStdoutStderrCompleter:
			return fmt.Errorf("CodeCrafters Internal Error: patching no_stdout_stderr_completer failed: %w", err)
		case secretPatchedEnvContextCompleter:
			return fmt.Errorf("CodeCrafters Internal Error: patching env_context_completer failed: %w", err)
		}
	}

	return reSignExecutableDarwinARM64(outputPath)
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
