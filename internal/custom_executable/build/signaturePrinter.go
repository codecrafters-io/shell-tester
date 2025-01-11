package custom_executable

import (
	"fmt"
	"os"
	"strings"
)

func addSecretCodeToExecutable(filePath, randomString string) error {
	LENGTH := 10
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: read file failed: %w", err)
	}
	placeholderIndex := strings.Index(string(data), "<<RANDOM>>")
	if placeholderIndex == -1 {
		return fmt.Errorf("CodeCrafters Internal Error: <<RANDOM>> not found in file")
	}
	bytes := copy(data[placeholderIndex:placeholderIndex+LENGTH], []byte(randomString))
	if bytes != LENGTH {
		return fmt.Errorf("CodeCrafters Internal Error: copy failed")
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: write file failed: %w", err)
	}
	return nil
}

func CreateSignaturePrinterExecutable(randomString, outputPath string) error {
	// Our executable contains a placeholder for the random string
	// The placeholder is <<<RANDOM>>>
	// We will replace the placeholder with the random string
	// The random string HAS to be the same length as the placeholder
	if len(randomString) != 10 {
		return fmt.Errorf("CodeCrafters Internal Error: randomString length must be 10")
	}

	// Copy the base executable from archive location to user's executable path
	err := copyExecutable("signature_printer", outputPath)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: copying executable failed: %w", err)
	}

	// Replace the placeholder with the random string
	// We can run the executable now, it will work as expected
	err = addSecretCodeToExecutable(outputPath, randomString)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: adding secret code to executable failed: %w", err)
	}

	return nil
}
