package custom_executable

import "fmt"

func CreateHeadExecutable(outputPath string) error {
	err := createExecutableForOSAndArch("head", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: creating executable %s failed: %w", "head", err)
	}
	return nil
}
