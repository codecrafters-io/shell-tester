package custom_executable

import "fmt"

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

	if err := reSignExecutableDarwinArm64(outputPath); err != nil {
		return err
	}

	return nil
}
