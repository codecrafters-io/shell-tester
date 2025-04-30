package custom_executable

import "fmt"

func CreateGrepExecutable(outputPath string) error {
	err := createExecutableForOSAndArch("grep", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: creating executable %s failed: %w", "grep", err)
	}
	return nil
}
