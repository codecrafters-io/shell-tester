package custom_executable

import "fmt"

func CreateYesExecutable(outputPath string) error {
	err := createExecutableForOSAndArch("yes", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: creating executable %s failed: %w", "yes", err)
	}
	return nil
}
