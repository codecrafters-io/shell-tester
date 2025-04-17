package custom_executable

import "fmt"

func CreateWcExecutable(outputPath string) error {
	err := createExecutableForOSAndArch("wc", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: creating executable %s failed: %w", "wc", err)
	}
	return nil
}
