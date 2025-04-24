package custom_executable

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func addSecretCodeToExecutable(filePath, randomString string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: read file failed: %w", err)
	}
	newData := strings.ReplaceAll(string(data), "<<RANDOM>>", randomString)
	if err := os.WriteFile(filePath, []byte(newData), 0644); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: write file failed: %w", err)
	}
	return nil
}

func CreateSignaturePrinterExecutable(randomString, outputPath string) error {
	// Our executable contains a placeholder for the random string
	// The placeholder is <<<RANDOM>>>
	// We will replace the placeholder with the random string
	// The random string HAS to be the same length as the placeholder
	if len(randomString) != 10 {
		return fmt.Errorf("CodeCrafters Internal Error: randomString length must be 10")
	}

	executableName := "signature_printer"
	// Copy the base executable from archive location to user's executable path
	err := createExecutableForOSAndArch(executableName, outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: copying executable failed: %w", err)
	}

	// Replace the placeholder with the random string
	// We can run the executable now, it will work as expected
	err = addSecretCodeToExecutable(outputPath, randomString)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: adding secret code to executable failed: %w", err)
	}

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		// Remove the signature from the executable
		command := fmt.Sprintf("codesign --remove-signature %s", outputPath)
		err = exec.Command("sh", "-c", command).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters Internal Error: removing signature from executable failed: %w", err)
		}

		// Sign the executable
		command = fmt.Sprintf("codesign -s - %s", outputPath)
		err = exec.Command("sh", "-c", command).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters Internal Error: signing executable failed: %w", err)
		}

		// Verify the signature
		command = fmt.Sprintf("codesign -vv %s", outputPath)
		err = exec.Command("sh", "-c", command).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters Internal Error: verifying signature failed: %w", err)
		}

		// Print the signature
		command = fmt.Sprintf("codesign -d --verbose=2 %s", outputPath)
		err = exec.Command("sh", "-c", command).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters Internal Error: printing signature failed: %w", err)
		}
	}

	return nil
}
