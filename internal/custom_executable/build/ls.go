package custom_executable

import "fmt"

func CreateLsExecutable(outputPath string) error {
	err := createExecutableForOSAndArch("ls", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: creating executable %s failed: %w", "ls", err)
	}
	return nil
}
