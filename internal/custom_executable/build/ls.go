package custom_executable

import "fmt"

func CreateLsExecutable(outputPath string) error {
	fileName := fetchCustomExecutableForOSAndArch("ls")

	// Copy the base executable from archive location to user's executable path
	err := copyExecutable(fileName, outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: copying executable %s failed: %w", fileName, err)
	}

	return nil
}
