package custom_executable

import "fmt"

func CreateCatExecutable(outputPath string) error {
	err := createExecutableForOSAndArch("cat", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: creating executable %s failed: %w", "cat", err)
	}

	return nil
}
