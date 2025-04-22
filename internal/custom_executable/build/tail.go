package custom_executable

import "fmt"

func CreateTailExecutable(outputPath string) error {
	err := createExecutableForOSAndArch("tail", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: creating executable %s failed: %w", "tail", err)
	}
	return nil
}
