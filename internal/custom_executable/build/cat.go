package custom_executable

import "fmt"

func CreateCatExecutable(outputPath string) error {
	fileName := fetchCustomExecutableForOSAndArch("cat")

	// Copy the base executable from archive location to user's executable path
	err := copyExecutable(fileName, outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: copying executable %s failed: %w", fileName, err)
	}

	return nil
}
