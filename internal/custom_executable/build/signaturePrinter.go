package custom_executable

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

	// We are okay with keeping this here
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		// Remove the signature from the executable
		err = exec.Command("codesign", "--remove-signature", outputPath).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters Internal Error: removing signature from executable failed: %w", err)
		}

		// Sign the executable
		err = exec.Command("codesign", "-s", "-", outputPath).Run()
		if err != nil {
			return fmt.Errorf("CodeCrafters Internal Error: signing executable failed: %w", err)
		}
	}

	return nil
}
